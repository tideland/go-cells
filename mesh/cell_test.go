// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"tideland.dev/go/asserts/verify"
)

// waitClosed waits for a channel to close or times out.
func waitClosed(t *testing.T, ch chan any, timeout time.Duration, msgAndArgs ...any) {
	t.Helper()
	select {
	case <-ch:
		// Channel closed successfully
	case <-time.After(timeout):
		msg := "channel not closed within timeout"
		if len(msgAndArgs) > 0 {
			if s, ok := msgAndArgs[0].(string); ok {
				msg = s
			}
		}
		t.Fatalf("%s", msg)
	}
}

// TestCellSimple provides a simple processing of some
// events.
func TestCellSimple(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	sigc := make(chan any, 1)
	collector := func(cell Cell, evt *Event, out Emitter) error {
		close(sigc)
		return nil
	}
	tbCollector := NewRequestBehavior(collector)
	cCollector := newCell(ctx, "collector", meshStub{}, tbCollector, drop)

	cCollector.Receive("one")

	waitClosed(t, sigc, time.Second)

	cancel()
}

// TestCellChain provides a chained processing of some
// events.
func TestCellChain(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	topics := []string{}
	sigc := make(chan any)
	upcaser := func(cell Cell, evt *Event, out Emitter) error {
		upperTopic := strings.ToUpper(evt.Topic())
		out.Emit(upperTopic)
		return nil
	}
	tbUpcaser := NewRequestBehavior(upcaser)
	cUpcaser := newCell(ctx, "upcaser", meshStub{}, tbUpcaser, drop)
	collector := func(cell Cell, evt *Event, out Emitter) error {
		topics = append(topics, evt.Topic())
		if len(topics) == 3 {
			close(sigc)
		}
		return nil
	}
	tbCollector := NewRequestBehavior(collector)
	cCollector := newCell(ctx, "collector", meshStub{}, tbCollector, drop)
	cCollector.SubscribeTo(cUpcaser)

	cUpcaser.Receive("one")
	cUpcaser.Receive("two")
	cUpcaser.Receive("three")

	waitClosed(t, sigc, time.Second)
	verify.Length(t, topics, 3)
	verify.Equal(t, strings.Join(topics, " "), "ONE TWO THREE")

	cCollector.UnsubscribeFrom(cUpcaser)

	cUpcaser.Receive("FOUR")
	cUpcaser.Receive("FIVE")
	cUpcaser.Receive("SIX")

	verify.Length(t, topics, 3)
	verify.Equal(t, strings.Join(topics, " "), "ONE TWO THREE")

	cancel()
}

// TestCellAutoUnsubscribe verifies the automatic unsubscription
// and information.
func TestCellAutoUnsubscribe(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	failed := []*Event{}
	collected := []*Event{}
	sigc := make(chan any)
	forwarder := func(cell Cell, evt *Event, out Emitter) error {
		return out.EmitEvent(evt)
	}
	cForwarderA := newCell(ctx, "forwarderA", meshStub{}, NewRequestBehavior(forwarder), drop)
	cForwarderB := newCell(ctx, "forwarderB", meshStub{}, NewRequestBehavior(forwarder), drop)
	failer := func(cell Cell, evt *Event, out Emitter) error {
		failed = append(failed, evt)
		if len(failed) == 3 {
			return errors.New("done")
		}
		return out.EmitEvent(evt)
	}
	cFailer := newCell(ctx, "failer", meshStub{}, NewRequestBehavior(failer), drop)
	cFailer.SubscribeTo(cForwarderA)
	cFailer.SubscribeTo(cForwarderB)
	collector := func(cell Cell, evt *Event, out Emitter) error {
		collected = append(collected, evt)
		if len(collected) == 3 {
			close(sigc)
		}
		return nil
	}
	cCollector := newCell(ctx, "collector", meshStub{}, NewRequestBehavior(collector), drop)
	cCollector.SubscribeTo(cFailer)

	cForwarderA.Receive("foo")
	cForwarderB.Receive("bar")
	cForwarderA.Receive("baz")

	waitClosed(t, sigc, time.Second)

	cForwarderA.Receive("dont-care")
	cForwarderB.Receive("dont-care")

	foundc := make(chan any)

	for _, evt := range collected {
		if evt.Topic() == TopicError {
			var errpl PayloadCellError
			err := evt.Payload(&errpl)
			verify.NoError(t, err)
			verify.Equal(t, errpl.CellName, "failer")
			verify.Equal(t, errpl.Error, "done")
			close(foundc)
			break
		}
	}

	waitClosed(t, foundc, time.Second, "error not found")
	cancel()
}

// meshStub simulates the mesh for the cells.
type meshStub struct{}

func (ms meshStub) Go(name string, b Behavior) error {
	return nil
}

func (ms meshStub) Subscribe(fromName, toName string) error {
	return nil
}

func (ms meshStub) Unsubscribe(toName, fromName string) error {
	return nil
}

func (ms meshStub) Emit(name, topic string, payloads ...any) error {
	return nil
}

func (ms meshStub) EmitEvent(name string, evt *Event) error {
	return nil
}

func (ms meshStub) Emitter(name string) (Emitter, error) {
	return nil, nil
}

// drop simulates the callback to notify the
// mesh of the termination of a cell.
var drop = func() {}
