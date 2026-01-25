// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"context"

	"tideland.dev/go/cells/mesh/internal"
)

// New creates new Mesh instance.
func New(ctx context.Context) Mesh {
	return internal.NewMesh(ctx)
}
