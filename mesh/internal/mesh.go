// Tideland Go Cells - Mesh - Internal
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package internal // import "tideland.dev/go/cells/mesh/internal"

//--------------------
// IMPORT
//--------------------

import (
	"context"
	"fmt"
	"sync"
)

//--------------------
// MESH
//--------------------

// Mesh manages a closed network of cells. It implements
// the mesh.Mesh interface.
type meshImpl struct {
	mu       sync.RWMutex
	ctx      context.Context
	cells    map[string]*CellImpl
	emitters map[string]*emitterImpl
}

// NewMesh creates new Mesh instance.
func NewMesh(ctx context.Context) *meshImpl {
	m := &meshImpl{
		ctx:      ctx,
		cells:    make(map[string]*CellImpl),
		emitters: make(map[string]*emitterImpl),
	}
	return m
}

// Go implements mesh.Mesh.
func (m *meshImpl) Go(name string, b Behavior) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.cells[name] != nil {
		return fmt.Errorf("cell name '%s' already used", name)
	}
	m.cells[name] = NewCell(m.ctx, name, m, b, func() {
		// Callback for cell to unregister.
		m.mu.Lock()
		defer m.mu.Unlock()
		delete(m.cells, name)
		delete(m.emitters, name)
	})
	return nil
}

// Subscribe implements mesh.Mesh.
func (m *meshImpl) Subscribe(emitterName, receptorName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitterCell := m.cells[emitterName]
	receptorCell := m.cells[receptorName]
	if emitterCell == nil {
		return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
	}
	if receptorCell == nil {
		return fmt.Errorf("receptor cell '%s' does not exist", receptorName)
	}
	receptorCell.SubscribeTo(emitterCell)
	return nil
}

// Unsubscribe implements mesh.Mesh.
func (m *meshImpl) Unsubscribe(emitterName, receptorName string) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitterCell := m.cells[emitterName]
	receptorCell := m.cells[receptorName]
	if emitterCell == nil {
		return fmt.Errorf("emitter cell '%s' does not exist", emitterName)
	}
	if receptorCell == nil {
		return fmt.Errorf("receptor cell '%s' does not exist", receptorName)
	}
	receptorCell.UnsubscribeFrom(emitterCell)
	return nil
}

// Emit implements Mesh.
func (m *meshImpl) Emit(name, topic string, payloads ...any) error {
	evt, err := NewEvent(topic, payloads...)
	if err != nil {
		return err
	}
	return m.EmitEvent(name, evt)
}

// EmitEvent implements Mesh.
func (m *meshImpl) EmitEvent(name string, evt *Event) error {
	m.mu.RLock()
	defer m.mu.RUnlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return fmt.Errorf("cell '%s' does not exist", name)
	}
	evt.initEmitters()
	return emitCell.ReceiveEvent(evt)
}

// Emitter implements mesh.Mesh.
func (m *meshImpl) Emitter(name string) (Emitter, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	emitCell := m.cells[name]
	if emitCell == nil {
		return nil, fmt.Errorf("cell '%s' does not exist", name)
	}
	namedEmitter := m.emitters[name]
	if namedEmitter == nil {
		namedEmitter = NewEmitter(emitCell)
		m.emitters[name] = namedEmitter
	}
	return namedEmitter, nil
}

// EOF
