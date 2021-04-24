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
	"errors"
	"testing"
	"time"

	"tideland.dev/go/audit/asserts"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TESTS
//--------------------

// TestTestbed verifies the working of the testbed for behavior tests.
func TestTestbed(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	forwarder := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				out.EmitEvent(evt)
			}
		}
	}
	behavior := mesh.BehaviorFunc(forwarder)
	eval := func(tbctx *mesh.TestbedContext, evt *mesh.Event) error {
		tbctx.EventSink().Push(evt)
		if tbctx.EventSink().Len() == 3 {
			// Done.
			tbctx.SetSuccess()
		}
		return nil
	}

	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("one")
		out.Emit("two")
		out.Emit("three")
	}, time.Second)
	assert.NoError(err)
	// assert.Equal(count, 3)
}

// TestTestbedMesh verifies the Mesh stubbing of the testbed.
func TestTestbedMesh(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	mesher := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				switch evt.Topic() {
				case "go":
					out.EmitEvent(evt)
					err := cell.Mesh().Go("anything", nil)
					assert.ErrorContains(err, "cell name 'anything' already used")
				case "subscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Subscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "unsubscribe":
					out.EmitEvent(evt)
					err := cell.Mesh().Unsubscribe("anything", "anything-else")
					assert.ErrorContains(err, "emitter cell 'anything' does not exist")
				case "emit":
					out.EmitEvent(evt)
					err := cell.Mesh().EmitEvent("anything", evt)
					assert.ErrorContains(err, "cell 'anything' does not exist")
				case "emitter":
					out.EmitEvent(evt)
					emtr, err := cell.Mesh().Emitter("anything")
					assert.ErrorContains(err, "cell 'anything' does not exist")
					assert.Nil(emtr)
				case "done":
					out.EmitEvent(evt)
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(mesher)
	eval := func(tbctx *mesh.TestbedContext, evt *mesh.Event) error {
		tbctx.EventSink().Push(evt)
		if tbctx.EventSink().Len() == 6 {
			tbctx.SetSuccess()
		}
		return nil
	}

	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("go")
		out.Emit("subscribe")
		out.Emit("unsubscribe")
		out.Emit("emit")
		out.Emit("emitter")
		out.Emit("done")
	}, time.Second)
	assert.NoError(err)
}

// TestTestbedError verifies the error handling of the Testbed.
func TestTestbedError(t *testing.T) {
	assert := asserts.NewTesting(t, asserts.FailStop)
	failer := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				switch evt.Topic() {
				case "go":
					out.Emit("ok")
				case "done":
					out.Emit("done")
				default:
					// Fail when topic is unknown.
					out.Emit("fail")
				}
			}
		}
	}
	behavior := mesh.BehaviorFunc(failer)
	eval := func(tbctx *mesh.TestbedContext, evt *mesh.Event) error {
		switch evt.Topic() {
		case "done":
			tbctx.SetSuccess()
			return nil
		case "fail":
			return errors.New("failure")
		default:
			return nil
		}
	}

	tb := mesh.NewTestbed(behavior, eval)
	err := tb.Go(func(out mesh.Emitter) {
		out.Emit("go")
		out.Emit("go")
		out.Emit("go")
		out.Emit("dunno!")
		out.Emit("go")
		out.Emit("go")
		out.Emit("done")
	}, time.Second)
	assert.ErrorContains(err, "failure")
}

// EOF
