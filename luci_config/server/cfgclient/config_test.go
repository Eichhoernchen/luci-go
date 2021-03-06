// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package cfgclient

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/luci/luci-go/common/config"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/luci_config/server/cfgclient/backend"

	"golang.org/x/net/context"

	. "github.com/smartystreets/goconvey/convey"
)

// testingBackend is a Backend implementation that ignores Authority.
type testingBackend struct {
	serviceURL *url.URL
	err        error
	items      []*backend.Item
	url        url.URL
	lastParams backend.Params
}

func (tb *testingBackend) ServiceURL(context.Context) url.URL { return *tb.serviceURL }

func (tb *testingBackend) Get(c context.Context, configSet, path string, p backend.Params) (*backend.Item, error) {
	tb.lastParams = p
	if err := tb.err; err != nil {
		return nil, tb.err
	}
	if len(tb.items) == 0 {
		return nil, ErrNoConfig
	}
	return tb.cloneItems()[0], nil
}

func (tb *testingBackend) GetAll(c context.Context, t backend.GetAllTarget, path string, p backend.Params) (
	[]*backend.Item, error) {

	tb.lastParams = p
	if err := tb.err; err != nil {
		return nil, tb.err
	}
	return tb.cloneItems(), nil
}

func (tb *testingBackend) ConfigSetURL(c context.Context, configSet string, p backend.Params) (url.URL, error) {
	tb.lastParams = p
	return tb.url, tb.err
}

func (tb *testingBackend) GetConfigInterface(c context.Context, a backend.Authority) config.Interface {
	panic("not supported")
}

func (tb *testingBackend) cloneItems() []*backend.Item {
	clones := make([]*backend.Item, len(tb.items))
	for i, it := range tb.items {
		clone := *it
		clones[i] = &clone
	}
	return clones
}

type errorMultiResolver struct {
	getErr func(i int) error
	out    *[]string
}

func (er *errorMultiResolver) PrepareMulti(size int) {
	*er.out = make([]string, size)
}

func (er *errorMultiResolver) ResolveItemAt(i int, it *backend.Item) error {
	if err := er.getErr(i); err != nil {
		return err
	}
	(*er.out)[i] = it.Content
	return nil
}

func TestConfig(t *testing.T) {
	t.Parallel()

	Convey(`A testing backend`, t, func() {
		c := context.Background()

		var err error
		tb := testingBackend{
			items: []*backend.Item{
				{Meta: backend.Meta{"projects/foo", "path", "####", "v1"}, Content: "foo"},
				{Meta: backend.Meta{"projects/bar", "path", "####", "v1"}, Content: "bar"},
			},
		}
		transformItems := func(fn func(int, *backend.Item)) {
			for i, itm := range tb.items {
				fn(i, itm)
			}
		}
		tb.serviceURL, err = url.Parse("http://example.com/config")
		if err != nil {
			panic(err)
		}
		c = backend.WithBackend(c, &tb)

		Convey(`Can resolve its service URL.`, func() {
			So(ServiceURL(c), ShouldResemble, *tb.serviceURL)
		})

		// Caching / type test cases.
		Convey(`Test byte resolver`, func() {
			transformItems(func(i int, ci *backend.Item) {
				ci.Content = string([]byte{byte(i)})
			})

			Convey(`Single`, func() {
				var val []byte
				So(Get(c, AsService, "", "", Bytes(&val), nil), ShouldBeNil)
				So(val, ShouldResemble, []byte{0})
			})

			Convey(`Multi`, func() {
				var (
					val  [][]byte
					meta []*Meta
				)
				So(Projects(c, AsService, "", BytesSlice(&val), &meta), ShouldBeNil)
				So(val, ShouldResemble, [][]byte{{0}, {1}})
				So(meta, ShouldResemble, []*Meta{
					{"projects/foo", "path", "####", "v1"},
					{"projects/bar", "path", "####", "v1"},
				})
			})
		})

		Convey(`Test string resolver`, func() {
			transformItems(func(i int, ci *backend.Item) {
				ci.Content = fmt.Sprintf("[%d]", i)
			})

			Convey(`Single`, func() {
				var val string
				So(Get(c, AsService, "", "", String(&val), nil), ShouldBeNil)
				So(val, ShouldEqual, "[0]")
			})

			Convey(`Multi`, func() {
				var (
					val  []string
					meta []*Meta
				)
				So(Projects(c, AsService, "", StringSlice(&val), &meta), ShouldBeNil)
				So(val, ShouldResemble, []string{"[0]", "[1]"})
				So(meta, ShouldResemble, []*Meta{
					{"projects/foo", "path", "####", "v1"},
					{"projects/bar", "path", "####", "v1"},
				})
			})
		})

		Convey(`Test nil resolver`, func() {
			transformItems(func(i int, ci *backend.Item) {
				ci.Content = fmt.Sprintf("[%d]", i)
			})

			Convey(`Single`, func() {
				var meta Meta
				So(Get(c, AsService, "", "", nil, &meta), ShouldBeNil)
				So(meta, ShouldResemble, Meta{"projects/foo", "path", "####", "v1"})
				So(tb.lastParams.Content, ShouldBeFalse)
			})

			Convey(`Multi`, func() {
				var meta []*Meta
				So(Projects(c, AsService, "", nil, &meta), ShouldBeNil)
				So(meta, ShouldResemble, []*Meta{
					{"projects/foo", "path", "####", "v1"},
					{"projects/bar", "path", "####", "v1"},
				})
				So(tb.lastParams.Content, ShouldBeFalse)
			})
		})

		Convey(`Projects with some errors returns a MultiError.`, func() {
			testErr := errors.New("test error")
			var (
				val  []string
				meta []*Meta
			)
			er := errorMultiResolver{
				out: &val,
				getErr: func(i int) error {
					if i == 1 {
						return testErr
					}
					return nil
				},
			}

			err := Projects(c, AsService, "", &er, &meta)

			// We have a MultiError with an error in position 1.
			So(err, ShouldResemble, errors.MultiError{nil, testErr})

			// One of our configs should have resolved.
			So(val, ShouldResemble, []string{"foo", ""})

			// Meta still works.
			So(meta, ShouldResemble, []*Meta{
				{"projects/foo", "path", "####", "v1"},
				{"projects/bar", "path", "####", "v1"},
			})
		})

		Convey(`Test ConfigSetURL`, func() {
			u, err := url.Parse("https://example.com/foo/bar")
			if err != nil {
				panic(err)
			}

			tb.url = *u
			uv, err := GetConfigSetURL(c, AsService, "")
			So(err, ShouldBeNil)
			So(uv, ShouldResemble, url.URL{Scheme: "https", Host: "example.com", Path: "/foo/bar"})
		})
	})
}
