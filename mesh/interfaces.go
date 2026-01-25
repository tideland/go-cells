// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"tideland.dev/go/cells/mesh/internal"
)

// RE-EXPORT INTERFACES

// Mesh describes the interface to a mesh of a cell from the
// perspective of a behavior.
type Mesh = internal.Mesh

// Cell describes the interface to a cell from the perspective
// of a behavior.
type Cell = internal.Cell

// Behavior describes what cell implementations must understand.
type Behavior = internal.Behavior

// Receptor defines the interface to receive events.
type Receptor = internal.Receptor

// Emitter defines the interface for emitting events to one
// or more cells.
type Emitter = internal.Emitter
