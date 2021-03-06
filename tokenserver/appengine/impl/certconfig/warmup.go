// Copyright 2017 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package certconfig

import (
	"golang.org/x/net/context"

	"github.com/luci/luci-go/server/warmup"
)

func init() {
	warmup.Register("tokenserver/appengine/impl/certconfig", func(c context.Context) error {
		// Warmup ID -> CN mapping. Ignore result (it may be "no such CA"), this is
		// fine, the mapping will be preloaded.
		GetCAByUniqueID(c, 0)
		return nil
	})
}
