// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh/internal"
)

//--------------------
// TOPIC CONSTANTS
//--------------------

// Standard event topics.
const (
	TopicTerminated        = internal.TopicTerminated
	TopicError             = internal.TopicError
	TopicTestbedDone       = internal.TopicTestbedDone
	TopicTestbedTerminated = internal.TopicTestbedTerminated
	TopicTestbedError      = internal.TopicTestbedError
)

//--------------------
// PAYLOAD TYPES
//--------------------

// PayloadTermination describes the normal termination of a cell.
type PayloadTermination = internal.PayloadTermination

// PayloadCellError describes the abnormal termination of a cell.
type PayloadCellError = internal.PayloadCellError

// EOF
