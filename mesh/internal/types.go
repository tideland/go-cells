// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

// Standard event topics.
const (
	TopicTerminated        = "terminated"
	TopicError             = "error"
	TopicTestbedDone       = "testbed-done"
	TopicTestbedTerminated = "testbed-terminated"
	TopicTestbedError      = "testbed-error"
)

// PayloadTermination describes the normal termination of a cell.
type PayloadTermination struct {
	CellName string `json:"cellName"`
}

// PayloadCellError describes the abnormal termination of a cell.
type PayloadCellError struct {
	CellName string `json:"cellName"`
	Error    string `json:"error"`
}
