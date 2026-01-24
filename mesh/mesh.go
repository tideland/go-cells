// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORT
//--------------------

import (
	"context"

	"tideland.dev/go/cells/mesh/internal"
)

//--------------------
// MESH
//--------------------

// New creates new Mesh instance.
func New(ctx context.Context) Mesh {
	return internal.NewMesh(ctx)
}

// EOF
