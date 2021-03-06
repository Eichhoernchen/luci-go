// Copyright 2015 The LUCI Authors. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

// Package memcache implements a caching config client backend backed by
// AppEngine's memcache service.
package memcache

import (
	"encoding/hex"
	"time"

	"github.com/luci/luci-go/common/errors"
	log "github.com/luci/luci-go/common/logging"
	"github.com/luci/luci-go/luci_config/server/cfgclient/backend"
	"github.com/luci/luci-go/luci_config/server/cfgclient/backend/caching"

	mc "github.com/luci/gae/service/memcache"

	"golang.org/x/net/context"
)

const (
	memCacheSchema  = "v1"
	maxMemCacheSize = 1024 * 1024 // 1MB
)

// Backend wraps a backend.B instance with a memcache-backed caching layer whose
// entries expire after exp.
func Backend(b backend.B, exp time.Duration) backend.B {
	return &caching.Backend{
		B: b,
		CacheGet: func(c context.Context, key caching.Key, l caching.Loader) (*caching.Value, error) {
			if key.Authority != backend.AsService {
				return l(c, key, nil)
			}

			// Is the item already cached?
			k := memcacheKey(&key)
			mci, err := mc.GetKey(c, k)
			switch err {
			case nil:
				// Value was cached, successfully retrieved.
				v, err := caching.DecodeValue(mci.Value())
				if err != nil {
					return nil, errors.Annotate(err).Reason("failed to decode cache value from %(key)q").
						D("key", k).Err()
				}
				return v, nil

			case mc.ErrCacheMiss:
				// Value was not cached. Load from Loader and cache.
				v, err := l(c, key, nil)
				if err != nil {
					return nil, err
				}

				// Attempt to cache the value. If this fails, we'll log a warning and
				// move on.
				err = func() error {
					d, err := v.Encode()
					if err != nil {
						return errors.Annotate(err).Reason("failed to encode value").Err()
					}

					if len(d) > maxMemCacheSize {
						return errors.Reason("entry exceeds memcache size (%(size)d > %(max)d)").
							D("size", len(d)).D("max", maxMemCacheSize).Err()
					}

					item := mc.NewItem(c, k).SetValue(d).SetExpiration(exp)
					if err := mc.Set(c, item); err != nil {
						return errors.Annotate(err).Err()
					}
					return nil
				}()
				if err != nil {
					log.Fields{
						log.ErrorKey: err,
						"key":        k,
					}.Warningf(c, "Failed to cache config.")
				}

				// Return the loaded value.
				return v, nil

			default:
				// Unknown memcache error.
				log.Fields{
					log.ErrorKey: err,
					"key":        k,
				}.Warningf(c, "Failed to decode memcached config.")
				return l(c, key, nil)
			}
		},
	}
}

func memcacheKey(key *caching.Key) string { return hex.EncodeToString(key.ParamHash()) }
