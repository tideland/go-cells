// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package internal

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
