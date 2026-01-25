// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"context"

	"tideland.dev/go/cells/mesh/internal"
)

// TestCellI interface for testing.
type TestCellI interface {
	Cell
	Receptor
	Emitter
	SubscribeTo(ic TestCellI)
	UnsubscribeFrom(ic TestCellI)
	Receive(topic string, payload ...any) error
	ReceiveEvent(evt *Event) error
}

// testCellWrapper wraps internal cell for testing.
type testCellWrapper struct {
	internal.TestCellI
}

func (tcw *testCellWrapper) SubscribeTo(ic TestCellI) {
	// Unwrap to get the underlying CellImpl
	if wrapped, ok := ic.(*testCellWrapper); ok {
		// Get the concrete CellImpl from the internal.TestCellI
		// Since NewTestCell returns a *CellImpl, we can type assert
		tcw.TestCellI.SubscribeTo(wrapped.TestCellI.(*internal.CellImpl))
	}
}

func (tcw *testCellWrapper) UnsubscribeFrom(ic TestCellI) {
	// Unwrap to get the underlying CellImpl
	if wrapped, ok := ic.(*testCellWrapper); ok {
		tcw.TestCellI.UnsubscribeFrom(wrapped.TestCellI.(*internal.CellImpl))
	}
}

// newCell is exported for testing purposes only.
func newCell(ctx context.Context, name string, m Mesh, b Behavior, drop func()) TestCellI {
	return &testCellWrapper{
		TestCellI: internal.NewTestCell(ctx, name, m, b, drop),
	}
}

// newStream is exported for testing purposes only.
func newStream() interface {
	Receptor
	Emitter
} {
	return internal.NewTestStream()
}
