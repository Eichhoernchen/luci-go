// Copyright 2016 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package signing

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"golang.org/x/net/context"

	"github.com/luci/luci-go/server/auth/identity"
	"github.com/luci/luci-go/server/auth/internal"

	. "github.com/luci/luci-go/common/testing/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchServiceInfo(t *testing.T) {
	Convey("Works", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(`{
				"app_id": "some-app-id",
				"app_runtime": "go",
				"app_runtime_version": "go1.5.1",
				"app_version": "1234-abcdef",
				"service_account_name": "some-app-id@appspot.gserviceaccount.com"
			}`))
		}))
		info, err := FetchServiceInfo(context.Background(), ts.URL)
		So(err, ShouldBeNil)
		So(info, ShouldResemble, &ServiceInfo{
			AppID:              "some-app-id",
			AppRuntime:         "go",
			AppRuntimeVersion:  "go1.5.1",
			AppVersion:         "1234-abcdef",
			ServiceAccountName: "some-app-id@appspot.gserviceaccount.com",
		})
	})

	Convey("Error", t, func() {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "fail", http.StatusInternalServerError)
		}))
		info, err := FetchServiceInfo(context.Background(), ts.URL)
		So(info, ShouldBeNil)
		So(err, ShouldNotBeNil)
	})
}

func TestFetchLUCIServiceIdentity(t *testing.T) {
	Convey("Works", t, func() {
		ctx := context.Background()
		ctx = internal.WithTestTransport(ctx, func(r *http.Request, body string) (code int, response string) {
			if r.URL.String() != "https://blah/auth/api/v1/server/info" {
				return 400, "Wrong URL"
			}
			return 200, `{"app_id": "blah-app-id"}`
		})

		id, err := FetchLUCIServiceIdentity(ctx, "https://blah")
		So(err, ShouldBeNil)
		So(id, ShouldEqual, identity.Identity("service:blah-app-id"))

		_, err = FetchLUCIServiceIdentity(ctx, "http://blah")
		So(err, ShouldErrLike, "not an https:// URL")

		_, err = FetchLUCIServiceIdentity(ctx, "https://blah/blah")
		So(err, ShouldErrLike, "not a root URL")

		_, err = FetchLUCIServiceIdentity(ctx, "https://")
		So(err, ShouldErrLike, "not a root URL")
	})
}
