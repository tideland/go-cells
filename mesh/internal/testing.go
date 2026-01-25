// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

import (
	"context"
)

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
