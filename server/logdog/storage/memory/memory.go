// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package memory

import (
	"errors"
	"sync"

	"github.com/luci/luci-go/common/logdog/types"
	"github.com/luci/luci-go/server/logdog/storage"
)

type logStream struct {
	logs        map[types.MessageIndex][]byte
	latestIndex types.MessageIndex
}

type rec struct {
	index types.MessageIndex
	data  []byte
}

// Storage is an implementation of the storage.Storage interface that stores
// data in memory.
//
// This is intended for testing, and not intended to be performant.
type Storage struct {
	// MaxGetCount, if not zero, is the maximum number of records to retrieve from
	// a single Get request.
	MaxGetCount int

	stateMu sync.Mutex
	streams map[types.StreamPath]*logStream
	closed  bool
}

var _ storage.Storage = (*Storage)(nil)

// Close implements storage.Storage.
func (s *Storage) Close() {
	s.run(func() error {
		s.closed = true
		return nil
	})
}

// Put implements storage.Storage.
func (s *Storage) Put(req *storage.PutRequest) error {
	return s.run(func() error {
		ls := s.getLogStreamLocked(req.Path, true)

		if _, ok := ls.logs[req.Index]; ok {
			return storage.ErrExists
		}

		ls.logs[req.Index] = []byte(req.Value)
		if req.Index > ls.latestIndex {
			ls.latestIndex = req.Index
		}
		return nil
	})
}

// Get implements storage.Storage.
func (s *Storage) Get(req *storage.GetRequest, cb storage.GetCallback) error {
	recs := []*rec(nil)
	err := s.run(func() error {
		ls := s.getLogStreamLocked(req.Path, false)
		if ls == nil {
			return storage.ErrDoesNotExist
		}

		limit := len(ls.logs)
		if req.Limit > 0 && req.Limit < limit {
			limit = req.Limit
		}
		if s.MaxGetCount > 0 && s.MaxGetCount < limit {
			limit = s.MaxGetCount
		}

		// Grab all records starting from our start index.
		for idx := req.Index; idx <= ls.latestIndex; idx++ {
			if le, ok := ls.logs[idx]; ok {
				recs = append(recs, &rec{
					index: idx,
					data:  le,
				})
			}

			if len(recs) >= limit {
				break
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// Punt all of the records upstream. We copy the data to prevent the
	// callback from accidentally mutating it. We reuse the data buffer to try
	// and catch errors when the callback retains the data.
	for _, r := range recs {
		dataCopy := make([]byte, len(r.data))
		copy(dataCopy, r.data)
		if !cb(r.index, dataCopy) {
			break
		}
	}

	return nil
}

// Tail implements storage.Storage.
func (s *Storage) Tail(p types.StreamPath) ([]byte, types.MessageIndex, error) {
	var r *rec

	// Find the latest log, then return it.
	err := s.run(func() error {
		ls := s.getLogStreamLocked(p, false)
		if ls == nil {
			return storage.ErrDoesNotExist
		}

		r = &rec{
			index: ls.latestIndex,
			data:  ls.logs[ls.latestIndex],
		}
		return nil
	})
	if err != nil {
		return nil, 0, err
	}
	return r.data, r.index, nil
}

// Purge implements storage.Storage.
func (s *Storage) Purge(p types.StreamPath) error {
	return s.run(func() error {
		if _, ok := s.streams[p]; !ok {
			return storage.ErrDoesNotExist
		}
		delete(s.streams, p)
		return nil
	})
}

func (s *Storage) run(f func() error) error {
	s.stateMu.Lock()
	defer s.stateMu.Unlock()

	if s.closed {
		return errors.New("storage is closed")
	}
	return f()
}

func (s *Storage) getLogStreamLocked(p types.StreamPath, create bool) *logStream {
	ls := s.streams[p]
	if ls == nil && create {
		ls = &logStream{
			logs:        map[types.MessageIndex][]byte{},
			latestIndex: -1,
		}

		if s.streams == nil {
			s.streams = map[types.StreamPath]*logStream{}
		}
		s.streams[p] = ls
	}

	return ls
}
