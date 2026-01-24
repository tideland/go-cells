// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// TOPIC CONSTANTS
//--------------------

// Standard event topics.
const (
	TopicTerminated        = "terminated"
	TopicError             = "error"
	TopicTestbedDone       = "testbed-done"
	TopicTestbedTerminated = "testbed-terminated"
	TopicTestbedError      = "testbed-error"
)

//--------------------
// PAYLOAD TYPES
//--------------------

// PayloadTermination describes the normal termination of a cell.
type PayloadTermination struct {
	CellName string `json:"cellName"`
}

// PayloadCellError describes the abnormal termination of a cell.
type PayloadCellError struct {
	CellName string `json:"cellName"`
	Error    string `json:"error"`
}

// EOF
