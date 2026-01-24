// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORTS
//--------------------


//--------------------
// EMITTER
//--------------------

// Emitter allows the continuous emitting of events to a cell
// without having to resolve the cell name each time.
type emitterImpl struct {
	cell *CellImpl
}

// NewEmitter creates a new emitter for the given cell.
func NewEmitter(cell *CellImpl) *emitterImpl {
	return &emitterImpl{
		cell: cell,
	}
}

// Emit implements Emitter.
func (e *emitterImpl) Emit(topic string, payloads ...any) error {
	return e.cell.Receive(topic, payloads...)
}

// EmitEvent implements Emitter.
func (e *emitterImpl) EmitEvent(evt *Event) error {
	return e.cell.ReceiveEvent(evt)
}

// EOF
