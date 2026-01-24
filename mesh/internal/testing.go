// Tideland Go Cells - Mesh - Internal - Testing
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
)

//--------------------
// TESTING EXPORTS
//--------------------

// TestCellI interface for testing.
type TestCellI interface {
	Cell
	Receptor
	Emitter
	SubscribeTo(ic *CellImpl)
	UnsubscribeFrom(ic *CellImpl)
	Receive(topic string, payload ...any) error
	ReceiveEvent(evt *Event) error
}

// NewTestCell is exported for testing purposes only.
func NewTestCell(ctx context.Context, name string, m Mesh, b Behavior, drop func()) TestCellI {
	return NewCell(ctx, name, m, b, drop)
}

// NewTestStream is exported for testing purposes only.
func NewTestStream() interface {
	Receptor
	Emitter
} {
	return NewStream()
}

// EOF
