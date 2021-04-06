// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestNewMesh verifies the simple creation of a mesh.
func TestNewMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	msh := mesh.New(ctx)

	assert.NotNil(msh)

	cancel()
}

// TestMeshGo verifies the starting of a cell via mesh.
func TestMeshGo(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		sigc <- cell.Name()
		return nil
	}
	msh := mesh.New(ctx)

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	assert.Wait(sigc, "testing", time.Second)

	cancel()
}

// TestMeshSubscriptions verifies the subscription and unsubscription
// of cells.
func TestMeshSubscriptions(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	forwardFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				out.EmitEvent(evt)
			}
		}
	}
	collectFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		topics := []string{}
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				topics = append(topics, evt.Topic())
				if len(topics) == 3 {
					sigc <- len(topics)
				}
			}
		}
	}
	msh := mesh.New(ctx)

	// Both cells do not exist.
	err := msh.Subscribe("forwarder", "collector-a")
	assert.ErrorContains(err, "cell 'forwarder' does not exist")

	msh.Go("forwarder", mesh.BehaviorFunc(forwardFunc))

	// One cell do not exist.
	err = msh.Subscribe("forwarder", "collector-a")
	assert.ErrorContains(err, "cell 'collector-a' does not exist")

	// Both cells exist.
	msh.Go("collector-a", mesh.BehaviorFunc(collectFunc))
	err = msh.Subscribe("forwarder", "collector-a")
	assert.NoError(err)

	msh.Emit("forwarder", "one")
	msh.Emit("forwarder", "two")
	msh.Emit("forwarder", "three")

	assert.Wait(sigc, 3, time.Second)

	// Unsubscribe one collector but subscribe a new one.
	err = msh.Unsubscribe("forwarder", "collector-a")
	assert.NoError(err)
	msh.Go("collector-b", mesh.BehaviorFunc(collectFunc))
	err = msh.Subscribe("forwarder", "collector-b")
	assert.NoError(err)

	msh.Emit("forwarder", "one")
	msh.Emit("forwarder", "two")
	msh.Emit("forwarder", "three")

	assert.Wait(sigc, 3, time.Second)

	// Unsubscribe not existing cell.
	err = msh.Unsubscribe("forwarder", "dont-exist")
	assert.ErrorContains(err, "cell 'dont-exist' does not exist")

	// Unsubscribe not subscribed cell.
	err = msh.Unsubscribe("forwarder", "collector-a")
	assert.NoError(err)

	// Unsubscribe subscribed cell.
	err = msh.Unsubscribe("forwarder", "collector-b")
	assert.NoError(err)

	cancel()
}

// TestMeshEmit verifies the emitting of events to one cell.
func TestMeshEmit(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		i := 0
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				i++
				if evt.Topic() == "get-i" {
					sigc <- i
				}
			}
		}
	}
	msh := mesh.New(ctx)
	err := msh.Emit("testing", "one")
	assert.ErrorContains(err, "cell 'testing' does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	err = msh.Emit("testing", "one")
	assert.NoError(err)

	msh.Emit("testing", "two")
	msh.Emit("testing", "three")
	msh.Emit("testing", "get-i")

	assert.Wait(sigc, 4, time.Second)

	cancel()
}

// TestMeshEmitter verifies the emitting of events to one cell using an emitter.
func TestMeshEmitter(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan interface{})
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		i := 0
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				i++
				if evt.Topic() == "get-i" {
					sigc <- i
				}
			}
		}
	}
	msh := mesh.New(ctx)
	emtr, err := msh.Emitter("testing")
	assert.ErrorContains(err, "cell 'testing' does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))
	emtr, err = msh.Emitter("testing")
	assert.NoError(err)

	emtr.Emit("one")
	emtr.Emit("two")
	emtr.Emit("three")
	emtr.Emit("get-i")

	assert.Wait(sigc, 4, time.Second)

	cancel()
}

// TestMeshStoppedCell verifies the handling of emittings to
// stopped cells.
func TestMeshStoppedCell(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	ctx, cancel := context.WithCancel(context.Background())
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		i := 0
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case <-in.Pull():
				i++
				if i >= 3 {
					return nil
				}
			}
		}
	}
	msh := mesh.New(ctx)
	msh.Go("countdown", mesh.BehaviorFunc(behaviorFunc))

	assert.NoError(msh.Emit("countdown", "one"))
	assert.NoError(msh.Emit("countdown", "two"))
	assert.NoError(msh.Emit("countdown", "three"))
	assert.ErrorContains(msh.Emit("countdown", "four"), "timeout")
	assert.ErrorContains(msh.Emit("countdown", "five"), "cell 'countdown' does not exist")

	cancel()
}

// EOF
