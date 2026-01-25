// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"tideland.dev/go/cells/mesh/internal"
)

// Standard event topics.
const (
	TopicTerminated        = internal.TopicTerminated
	TopicError             = internal.TopicError
	TopicTestbedDone       = internal.TopicTestbedDone
	TopicTestbedTerminated = internal.TopicTestbedTerminated
	TopicTestbedError      = internal.TopicTestbedError
)

// PayloadTermination describes the normal termination of a cell.
type PayloadTermination = internal.PayloadTermination

// PayloadCellError describes the abnormal termination of a cell.
type PayloadCellError = internal.PayloadCellError
