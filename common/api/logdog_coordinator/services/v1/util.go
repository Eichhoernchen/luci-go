// Copyright 2016 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

package logdog

// Complete returns true if the archive request expresses that the archived
// log stream was complete.
//
// A log stream is complete if every entry between zero and its terminal index
// is included.
func (ar *ArchiveStreamRequest) Complete() bool {
	tidx := ar.TerminalIndex
	if tidx < 0 {
		tidx = -1
	}
	return (ar.LogEntryCount == (tidx + 1))
}