// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// EMITTER
//--------------------

// emitter allows the continuous emitting of events to a cell
// without having to resolve the cell name each time.
type emitter struct {
	cell *cell
}

// Emit implements Emitter.
func (e *emitter) Emit(topic string, payloads ...any) error {
	return e.cell.receive(topic, payloads...)
}

// EmitEvent implements Emitter.
func (e *emitter) EmitEvent(evt *Event) error {
	return e.cell.receiveEvent(evt)
}

// EOF
