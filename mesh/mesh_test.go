// Tideland Go Cells - Mesh - Tests
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh_test // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"context"
	"reflect"
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// HELPERS
//--------------------

// wait waits for a channel to receive a specific value or times out.
func wait(t *testing.T, ch chan any, expected any, timeout time.Duration, msgAndArgs ...any) {
	t.Helper()
	select {
	case got := <-ch:
		if !reflect.DeepEqual(got, expected) {
			t.Fatalf("expected %v, got %v", expected, got)
		}
	case <-time.After(timeout):
		msg := "channel did not receive expected value within timeout"
		if len(msgAndArgs) > 0 {
			if s, ok := msgAndArgs[0].(string); ok {
				msg = s
			}
		}
		t.Fatalf("%s", msg)
	}
}

//--------------------
// TESTS
//--------------------

// TestNewMesh verifies the simple creation of a mesh.
func TestNewMesh(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	msh := mesh.New(ctx)

	verify.NotNil(t,msh)

	cancel()
}

// TestMeshGo verifies the starting of a cell via mesh.
func TestMeshGo(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
	behaviorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		sigc <- cell.Name()
		return nil
	}
	msh := mesh.New(ctx)

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	wait(t,sigc, "testing", time.Second)

	cancel()
}

// TestMeshSubscriptions verifies the subscription and unsubscription
// of cells.
func TestMeshSubscriptions(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
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
	verify.ErrorContains(t,err, "cell 'forwarder' does not exist")

	msh.Go("forwarder", mesh.BehaviorFunc(forwardFunc))

	// One cell do not exist.
	err = msh.Subscribe("forwarder", "collector-a")
	verify.ErrorContains(t,err, "cell 'collector-a' does not exist")

	// Both cells exist.
	msh.Go("collector-a", mesh.BehaviorFunc(collectFunc))
	err = msh.Subscribe("forwarder", "collector-a")
	verify.NoError(t,err)

	msh.Emit("forwarder", "one")
	msh.Emit("forwarder", "two")
	msh.Emit("forwarder", "three")

	wait(t,sigc, 3, time.Second)

	// Unsubscribe one collector but subscribe a new one.
	err = msh.Unsubscribe("forwarder", "collector-a")
	verify.NoError(t,err)
	msh.Go("collector-b", mesh.BehaviorFunc(collectFunc))
	err = msh.Subscribe("forwarder", "collector-b")
	verify.NoError(t,err)

	msh.Emit("forwarder", "one")
	msh.Emit("forwarder", "two")
	msh.Emit("forwarder", "three")

	wait(t,sigc, 3, time.Second)

	// Unsubscribe not existing cell.
	err = msh.Unsubscribe("forwarder", "dont-exist")
	verify.ErrorContains(t,err, "cell 'dont-exist' does not exist")

	// Unsubscribe not subscribed cell.
	err = msh.Unsubscribe("forwarder", "collector-a")
	verify.NoError(t,err)

	// Unsubscribe subscribed cell.
	err = msh.Unsubscribe("forwarder", "collector-b")
	verify.NoError(t,err)

	cancel()
}

// TestMeshEmit verifies the emitting of events to one cell.
func TestMeshEmit(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
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
	verify.ErrorContains(t,err, "cell 'testing' does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))

	err = msh.Emit("testing", "one")
	verify.NoError(t,err)

	msh.Emit("testing", "two")
	msh.Emit("testing", "three")
	msh.Emit("testing", "get-i")

	wait(t,sigc, 4, time.Second)

	cancel()
}

// TestMeshEmitter verifies the emitting of events to one cell using an emitter.
func TestMeshEmitter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
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
	verify.ErrorContains(t,err, "cell 'testing' does not exist")

	msh.Go("testing", mesh.BehaviorFunc(behaviorFunc))
	emtr, err = msh.Emitter("testing")
	verify.NoError(t,err)

	emtr.Emit("one")
	emtr.Emit("two")
	emtr.Emit("three")
	emtr.Emit("get-i")

	wait(t,sigc, 4, time.Second)

	cancel()
}

// TestMeshStoppedCell verifies the handling of emittings to
// stopped cells.
func TestMeshStoppedCell(t *testing.T) {
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

	verify.NoError(t,msh.Emit("countdown", "one"))
	verify.NoError(t,msh.Emit("countdown", "two"))
	verify.NoError(t,msh.Emit("countdown", "three"))
	verify.ErrorContains(t,msh.Emit("countdown", "four"), "timeout")
	verify.ErrorContains(t,msh.Emit("countdown", "five"), "cell 'countdown' does not exist")

	cancel()
}

// TestMeshEmitters verifies different emittings and re-emittings
// and the according emitting entries.
func TestMeshEnitters(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
	collectorFunc := func(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
		emitters := make(map[string]bool)
		for {
			select {
			case <-cell.Context().Done():
				return nil
			case evt := <-in.Pull():
				emitters[cell.Name()+" :: "+evt.Topic()+" :: "+evt.Emitters()] = true
				switch evt.Topic() {
				case "emit":
					out.Emit("emitted")
				case "re-emit":
					if cell.Name() != "third" {
						out.EmitEvent(evt)
					}
				case "done":
					if evt.HasPayload() {
						var pl map[string]bool
						if err := evt.Payload(&pl); err != nil {
							return err
						}
						for emitter := range pl {
							emitters[emitter] = true
						}
					}
					if cell.Name() == "third" {
						sigc <- emitters
					} else {
						out.Emit("done", emitters)
					}
				}
			}
		}
	}
	msh := mesh.New(ctx)
	msh.Go("first", mesh.BehaviorFunc(collectorFunc))
	msh.Go("second", mesh.BehaviorFunc(collectorFunc))
	msh.Go("third", mesh.BehaviorFunc(collectorFunc))
	msh.Subscribe("first", "second")
	msh.Subscribe("second", "third")

	msh.Emit("first", "anything")
	msh.Emit("first", "emit")
	msh.Emit("first", "re-emit")
	msh.Emit("first", "done")

	wait(t,sigc, map[string]bool{
		"first :: anything :: /":            true,
		"first :: emit :: /":                true,
		"first :: re-emit :: /":             true,
		"first :: done :: /":                true,
		"second :: emitted :: first":        true,
		"second :: re-emit :: /first":       true,
		"second :: done :: first":           true,
		"third :: re-emit :: /first/second": true,
		"third :: done :: second":           true,
	}, time.Second)

	cancel()
}

// EOF
