// Copyright 2017 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package localauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"

	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/clock/testclock"
	"github.com/luci/luci-go/common/errors"
	"github.com/luci/luci-go/lucictx"

	. "github.com/luci/luci-go/common/testing/assertions"
	. "github.com/smartystreets/goconvey/convey"
)

func TestServerLifecycle(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	Convey("Serve or close before init", t, func() {
		s := Server{}
		So(s.Serve(), ShouldErrLike, "not initialized")
		So(s.Close(), ShouldErrLike, "not initialized")
	})

	Convey("Double init", t, func() {
		s := Server{}
		defer s.Close()
		_, err := s.Initialize(ctx)
		So(err, ShouldBeNil)
		_, err = s.Initialize(ctx)
		So(err, ShouldErrLike, "already initialized")
	})

	Convey("Server after close", t, func() {
		s := Server{}
		_, err := s.Initialize(ctx)
		So(err, ShouldBeNil)
		So(s.Close(), ShouldBeNil)
		So(s.Serve(), ShouldErrLike, "already closed")
	})

	Convey("Close works", t, func() {
		serving := make(chan struct{})
		s := Server{
			testingServeHook: func() { close(serving) },
		}
		_, err := s.Initialize(ctx)
		So(err, ShouldBeNil)

		done := make(chan error)
		go func() {
			done <- s.Serve()
		}()

		<-serving // wait until really started

		// Stop it.
		So(s.Close(), ShouldBeNil)
		// Doing it second time is ok too.
		So(s.Close(), ShouldBeNil)

		serverErr := <-done
		So(serverErr, ShouldBeNil)
	})
}

func TestProtocol(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	ctx, _ = testclock.UseTime(ctx, testclock.TestRecentTimeLocal)

	Convey("With server", t, func(c C) {
		// Use channels to pass mocked requests/responses back and forth.
		requests := make(chan []string, 10000)
		responses := make(chan interface{}, 1)
		s := Server{
			TokenGenerator: func(ctx context.Context, scopes []string, lifetime time.Duration) (*oauth2.Token, error) {
				requests <- scopes
				var resp interface{}
				select {
				case resp = <-responses:
				default:
					c.Println("Unexpected token request")
					return nil, fmt.Errorf("Unexpected request")
				}
				switch resp := resp.(type) {
				case error:
					return nil, resp
				case *oauth2.Token:
					return resp, nil
				default:
					panic("unknown response")
				}
			},
		}
		p, err := s.Initialize(ctx)
		So(err, ShouldBeNil)

		done := make(chan struct{})
		go func() {
			s.Serve()
			close(done)
		}()
		defer func() {
			s.Close()
			<-done
		}()

		goodRequest := func() *http.Request {
			return prepReq(p, "/rpc/LuciLocalAuthService.GetOAuthToken", map[string]interface{}{
				"scopes": []string{"B", "A"},
				"secret": p.Secret,
			})
		}

		Convey("Happy path", func() {
			responses <- &oauth2.Token{
				AccessToken: "tok1",
				Expiry:      clock.Now(ctx).Add(30 * time.Minute),
			}
			So(call(goodRequest()), ShouldEqual, `HTTP 200 (json): {"access_token":"tok1","expiry":1454502906}`)
			So(<-requests, ShouldResemble, []string{"A", "B"})

			// application/json is also the default.
			req := goodRequest()
			req.Header.Del("Content-Type")
			responses <- &oauth2.Token{
				AccessToken: "tok2",
				Expiry:      clock.Now(ctx).Add(30 * time.Minute),
			}
			So(call(req), ShouldEqual, `HTTP 200 (json): {"access_token":"tok2","expiry":1454502906}`)
			So(<-requests, ShouldResemble, []string{"A", "B"})
		})

		Convey("Panic in token generator", func() {
			responses <- "omg, panic"
			So(call(goodRequest()), ShouldEqual, `HTTP 500: Internal Server Error. See logs.`)
		})

		Convey("Not POST", func() {
			req := goodRequest()
			req.Method = "PUT"
			So(call(req), ShouldEqual, `HTTP 405: Expecting POST`)
		})

		Convey("Bad URI", func() {
			req := goodRequest()
			req.URL.Path = "/zzz"
			So(call(req), ShouldEqual, `HTTP 404: Expecting /rpc/LuciLocalAuthService.<method>`)
		})

		Convey("Bad content type", func() {
			req := goodRequest()
			req.Header.Set("Content-Type", "bzzzz")
			So(call(req), ShouldEqual, `HTTP 400: Expecting 'application/json' Content-Type`)
		})

		Convey("Broken json", func() {
			req := goodRequest()

			body := `not a json`
			req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
			req.ContentLength = int64(len(body))

			So(call(req), ShouldEqual, `HTTP 400: Not JSON body - invalid character 'o' in literal null (expecting 'u')`)
		})

		Convey("Huge request", func() {
			req := goodRequest()

			body := strings.Repeat("z", 64*1024+1)
			req.Body = ioutil.NopCloser(bytes.NewBufferString(body))
			req.ContentLength = int64(len(body))

			So(call(req), ShouldEqual, `HTTP 400: Expecting 'Content-Length' header, <64Kb`)
		})

		Convey("Unknown RPC method", func() {
			req := prepReq(p, "/rpc/LuciLocalAuthService.UnknownMethod", map[string]interface{}{})
			So(call(req), ShouldEqual, `HTTP 404: Unknown RPC method "UnknownMethod"`)
		})

		Convey("No scopes", func() {
			req := prepReq(p, "/rpc/LuciLocalAuthService.GetOAuthToken", map[string]interface{}{
				"secret": p.Secret,
			})
			So(call(req), ShouldEqual, `HTTP 400: Field "scopes" is required.`)
		})

		Convey("No secret", func() {
			req := prepReq(p, "/rpc/LuciLocalAuthService.GetOAuthToken", map[string]interface{}{
				"scopes": []string{"B", "A"},
			})
			So(call(req), ShouldEqual, `HTTP 400: Field "secret" is required.`)
		})

		Convey("Bad secret", func() {
			req := prepReq(p, "/rpc/LuciLocalAuthService.GetOAuthToken", map[string]interface{}{
				"scopes": []string{"B", "A"},
				"secret": []byte{0, 1, 2, 3},
			})
			So(call(req), ShouldEqual, `HTTP 403: Invalid secret.`)
		})

		Convey("Token generator returns fatal error", func() {
			responses <- fmt.Errorf("fatal!!111")
			So(call(goodRequest()), ShouldEqual, `HTTP 200 (json): {"error_code":-1,"error_message":"fatal!!111"}`)
		})

		Convey("Token generator returns ErrorWithCode", func() {
			responses <- errWithCode{
				error: fmt.Errorf("with code"),
				code:  123,
			}
			So(call(goodRequest()), ShouldEqual, `HTTP 200 (json): {"error_code":123,"error_message":"with code"}`)
		})

		Convey("Token generator returns transient error", func() {
			responses <- errors.WrapTransient(fmt.Errorf("transient"))
			So(call(goodRequest()), ShouldEqual, `HTTP 500: Transient error - transient`)
		})
	})
}

type errWithCode struct {
	error
	code int
}

func (e errWithCode) Code() int {
	return e.code
}

func prepReq(p *lucictx.LocalAuth, uri string, body interface{}) *http.Request {
	var reader io.Reader
	isJSON := false
	if body != nil {
		blob, ok := body.([]byte)
		if !ok {
			var err error
			blob, err = json.Marshal(body)
			if err != nil {
				panic(err)
			}
			isJSON = true
		}
		reader = bytes.NewReader(blob)
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("http://127.0.0.1:%d%s", p.RPCPort, uri), reader)
	if err != nil {
		panic(err)
	}
	if isJSON {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func call(req *http.Request) interface{} {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	blob, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	tp := ""
	if resp.Header.Get("Content-Type") == "application/json; charset=utf-8" {
		tp = " (json)"
	}

	return fmt.Sprintf("HTTP %d%s: %s", resp.StatusCode, tp, strings.TrimSpace(string(blob)))
}
