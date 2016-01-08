// Copyright 2015 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package fetcher

import (
	"fmt"
	"io"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/luci/luci-go/common/clock"
	"github.com/luci/luci-go/common/logdog/protocol"
	"github.com/luci/luci-go/common/logdog/types"
	log "github.com/luci/luci-go/common/logging"
	"golang.org/x/net/context"
)

const (
	// DefaultDelay is the default Delay value.
	DefaultDelay = 5 * time.Second

	// DefaultBufferBytes is the default number of bytes to buffer.
	DefaultBufferBytes = int64(1024 * 1024) // 1MB
)

// LogRequest is a structure used by the Fetcher to request a range of logs
// from its Source.
type LogRequest struct {
	// Index is the starting log index to request.
	Index types.MessageIndex
	// Count, if >0, is the maximum number of log entries to request.
	Count int
	// Bytes, if >0, is the maximum number of log bytes to request. At least one
	// log must be returned regardless of the byte limit.
	Bytes int64
}

// Source is the source of log stream and log information.
//
// The Source is resposible for handling retries, backoff, and transient errors.
// An error from the Source will shut down the Fetcher.
type Source interface {
	// LogEntries populates the supplied LogRequest with available sequential
	// log entries as available. This will block until at least one entry is
	// available.
	//
	// Upon success, the requested logs and terminal message index is returned. If
	// no terminal index is known, a value <0 will be returned.
	LogEntries(context.Context, *LogRequest) ([]*protocol.LogEntry, types.MessageIndex, error)
}

// Options is the set of configuration parameters for a Fetcher.
type Options struct {
	// Source is the log stream source.
	Source Source

	// Index is the starting stream index to retrieve.
	Index types.MessageIndex
	// Count is the total number of logs to retrieve. If zero, the full stream
	// will be fetched.
	Count int64

	// Count is the minimum amount of LogEntry records to buffer. If the buffered
	// amount dips below Count, more logs will be fetched.
	//
	// If zero, no count target will be applied.
	BufferCount int

	// BufferBytes is the target number of LogEntry bytes to buffer. If the
	// buffered amount dips below Bytes, more logs will be fetched.
	//
	// If zero, no byte target will be applied unless Count is also zero, in which
	// case DefaultBufferBytes byte constraint will be applied.
	BufferBytes int64

	// PrefetchFactor constrains the amount of additional data to fetch when
	// refilling the buffer. Effective Count and Bytes values are multiplied
	// by PrefetchFactor to determine the amount of logs to request.
	PrefetchFactor int

	// Delay is the amount of time to wait in between unsuccessful log requests.
	Delay time.Duration
}

// A Fetcher buffers LogEntry records by querying the Source for log data.
// It attmepts to maintain a steady stream of records by prefetching available
// records in advance of consumption.
//
// A Fetcher is not goroutine-safe.
type Fetcher struct {
	// o is the set of Options.
	o *Options

	// logC is the communication channel between the fetch goroutine and
	// the user-facing NextLogEntry. Logs are punted through this channel
	// one-by-one.
	logC chan *logResponse
	// fetchErr is the retained error state. If not nil, fetching has stopped and
	// all Fetcher methods will return this error.
	fetchErr error

	// sizeFunc is a function that calculates the byte size of a LogEntry
	// protobuf.
	//
	// If nil, proto.Size will be used. This is used for testing.
	sizeFunc func(proto.Message) int
}

// New instantiates a new Fetcher instance.
//
// The Fetcher can be cancelled by cancelling the supplied context.
func New(c context.Context, o Options) *Fetcher {
	if o.Delay <= 0 {
		o.Delay = DefaultDelay
	}
	if o.BufferBytes <= 0 && o.BufferCount <= 0 {
		o.BufferBytes = DefaultBufferBytes
	}
	if o.PrefetchFactor <= 1 {
		o.PrefetchFactor = 1
	}

	f := Fetcher{
		o:    &o,
		logC: make(chan *logResponse, o.Count),
	}
	go f.fetch(c)
	return &f
}

type logResponse struct {
	log   *protocol.LogEntry
	size  int64
	index types.MessageIndex
	err   error
}

// NextLogEntry returns the next buffered LogEntry, blocking until it becomes
// available.
//
// If the end of the log stream is encountered, NextLogEntry will return
// io.EOF.
//
// If the Fetcher is cancelled, a context.Canceled error will be returned.
func (f *Fetcher) NextLogEntry() (*protocol.LogEntry, error) {
	var le *protocol.LogEntry
	if f.fetchErr == nil {
		resp := <-f.logC
		if resp.err != nil {
			f.fetchErr = resp.err
		}
		if resp.log != nil {
			le = resp.log
		}
	}
	return le, f.fetchErr
}

type fetchRequest struct {
	index types.MessageIndex
	count int
	bytes int64
	respC chan *fetchResponse
}

type fetchResponse struct {
	logs []*protocol.LogEntry
	tidx types.MessageIndex
	err  error
}

// fetch is run in a separate goroutine. It handles buffering goroutine and
// punts logs and/or error state to NextLogEntry via logC.
func (f *Fetcher) fetch(c context.Context) {
	defer close(f.logC)

	lb := logBuffer{}
	nextFetchIndex := f.o.Index
	lastSentIndex := types.MessageIndex(-1)
	tidx := types.MessageIndex(-1)

	c, cancelFunc := context.WithCancel(c)
	bytes := int64(0)
	logFetchBaseC := make(chan *fetchResponse)
	var logFetchC chan *fetchResponse
	defer func() {
		// Reap our fetch goroutine, if still fetching.
		if logFetchC != nil {
			cancelFunc()
			<-logFetchC
		}
		close(logFetchBaseC)
	}()

	errOut := func(err error) {
		f.logC <- &logResponse{err: err}
	}

	count := int64(0)
	for tidx < 0 || lastSentIndex < tidx {
		if f.o.Count > 0 && count >= f.o.Count {
			log.Fields{
				"count": count,
			}.Debugf(c, "Exceeded maximum log count.")
			break
		}

		// If we have a buffered log, prepare it for sending.
		var sendC chan<- *logResponse
		var lr *logResponse
		if le := lb.current(); le != nil {
			lr = &logResponse{
				log:   le,
				size:  f.sizeOf(le),
				index: types.MessageIndex(le.StreamIndex),
			}
			sendC = f.logC
		}

		// If we're not currently fetching logs, and we are below our thresholds,
		// request a new batch of logs.
		if logFetchC == nil && (tidx < 0 || nextFetchIndex <= tidx) &&
			(lb.size() < f.o.BufferCount || bytes < f.o.BufferBytes) {
			logFetchC = logFetchBaseC
			req := fetchRequest{
				index: nextFetchIndex,
				count: (f.o.BufferCount * f.o.PrefetchFactor) - lb.size(),
				bytes: (f.o.BufferBytes * int64(f.o.PrefetchFactor)) - bytes,
				respC: logFetchC,
			}

			doFetch := true
			if f.o.Count > 0 {
				// Constrain our logs to the maximim that we'll actually need.
				remaining := f.o.Count - count - int64(lb.size())
				if req.count == 0 || remaining < int64(req.count) {
					req.count = int(remaining)
				}
				if req.count <= 0 {
					// We have enough buffered logs.
					doFetch = false
				}
			}
			if doFetch {
				log.Fields{
					"index": req.index,
					"count": req.count,
					"bytes": req.bytes,
				}.Debugf(c, "Buffering next round of logs...")
				go f.fetchLogs(c, &req)
			}
		}

		select {
		case <-c.Done():
			// Our Context has been cancelled. Terminate and propagate that error.
			log.WithError(c.Err()).Debugf(c, "Context has been cancelled.")
			errOut(c.Err())
			return

		case resp := <-logFetchC:
			logFetchC = nil
			if resp.err != nil {
				log.WithError(resp.err).Errorf(c, "Error fetching logs.")
				errOut(resp.err)
				return
			}

			// Account for each log and add it to our buffer.
			if len(resp.logs) > 0 {
				for _, le := range resp.logs {
					bytes += int64(f.sizeOf(le))
				}
				nextFetchIndex = types.MessageIndex(resp.logs[len(resp.logs)-1].StreamIndex) + 1
				lb.append(resp.logs...)
			}
			if resp.tidx >= 0 {
				tidx = resp.tidx
			}

		case sendC <- lr:
			bytes -= lr.size
			lastSentIndex = lr.index

			count++
			lb.next()
		}
	}

	// We've punted the full log stream.
	errOut(io.EOF)
}

func (f *Fetcher) fetchLogs(c context.Context, req *fetchRequest) {
	resp := fetchResponse{}
	defer func() {
		if r := recover(); r != nil {
			resp.err = fmt.Errorf("panic during fetch: %v", r)
		}
		req.respC <- &resp
	}()

	// Request the next chunk of logs from our Source.
	lreq := LogRequest{
		Index: req.index,
		Count: req.count,
		Bytes: req.bytes,
	}
	if lreq.Count < 0 {
		lreq.Count = 0
	}
	if lreq.Bytes < 0 {
		lreq.Bytes = 0
	}

	for {
		log.Fields{
			"index": req.index,
		}.Debugf(c, "Fetching logs.")
		resp.logs, resp.tidx, resp.err = f.o.Source.LogEntries(c, &lreq)
		switch {
		case resp.err != nil:
			log.Fields{
				log.ErrorKey: resp.err,
				"index":      req.index,
			}.Errorf(c, "Fetch returned error.")
			return

		case len(resp.logs) > 0:
			log.Fields{
				"index":         req.index,
				"count":         len(resp.logs),
				"terminalIndex": resp.tidx,
				"firstLog":      resp.logs[0].StreamIndex,
				"lastLog":       resp.logs[len(resp.logs)-1].StreamIndex,
			}.Debugf(c, "Fetched log entries.")
			return

		case resp.tidx >= 0 && req.index > resp.tidx:
			// Will never be satisfied, since we're requesting beyond the terminal
			// index.
			return

		default:
			// No logs this round. Sleep for more.
			log.Fields{
				"index": req.index,
				"delay": f.o.Delay,
			}.Infof(c, "No logs returned. Sleeping...")
			clock.Sleep(c, f.o.Delay)
		}
	}
}

func (f *Fetcher) sizeOf(le *protocol.LogEntry) int64 {
	sf := f.sizeFunc
	if sf == nil {
		sf = proto.Size
	}
	return int64(sf(le))
}
