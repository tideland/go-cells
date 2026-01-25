// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh_test

import (
	"encoding/json"
	"testing"

	"tideland.dev/go/asserts/verify"

	"tideland.dev/go/cells/mesh"
)

// TestEventSimple verifies events without payloads.
func TestEventSimple(t *testing.T) {
	evt, err := mesh.NewEvent("")
	verify.ErrorContains(t, err, "event needs topic")
	if evt != nil {
		t.Fatalf("expected nil event, got %v", evt)
	}

	evt, err = mesh.NewEvent("test")
	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "test")
	verify.False(t, evt.HasPayload())
}

// TestEventPayload verifies events with one or more payloads.
func TestEventPayload(t *testing.T) {
	payloadIn := []string{"a", "b", "c"}
	payloadOutA := []string{}
	evt, err := mesh.NewEvent("test", payloadIn)
	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "test")
	verify.True(t, evt.HasPayload())
	err = evt.Payload(&payloadOutA)
	verify.NoError(t, err)
	verify.Length(t, payloadOutA, 3)
	verify.SliceEqual(t, payloadOutA, payloadIn)

	payloadOutB := []int{}
	evt, err = mesh.NewEvent("test", 1, 2, 3, 4, 5)
	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "test")
	verify.True(t, evt.HasPayload())
	err = evt.Payload(&payloadOutB)
	verify.NoError(t, err)
	verify.Length(t, payloadOutB, 5)
	verify.SliceEqual(t, payloadOutB, []int{1, 2, 3, 4, 5})

	var payloadOutC string
	evt, err = mesh.NewEvent("test", "payload")
	verify.NoError(t, err)
	verify.Equal(t, evt.Topic(), "test")
	verify.True(t, evt.HasPayload())
	err = evt.Payload(&payloadOutC)
	verify.NoError(t, err)
	verify.Equal(t, payloadOutC, "payload")

}

// TestEventMarshaling verifies the event marshaling and unmarshaling.
func TestEventMarshaling(t *testing.T) {
	evtIn, err := mesh.NewEvent("test")
	verify.NoError(t, err)
	data, err := json.Marshal(evtIn)
	verify.NoError(t, err)

	evtOut, err := mesh.NewEvent("empty")
	verify.NoError(t, err)
	err = json.Unmarshal(data, &evtOut)
	verify.NoError(t, err)
	verify.DeepEqual(t, evtOut, evtIn)

	plEvtA, err := mesh.NewEvent("payload-a")
	verify.NoError(t, err)
	plEvtB, err := mesh.NewEvent("payload-b")
	verify.NoError(t, err)
	plEvtC, err := mesh.NewEvent("payload-c")
	verify.NoError(t, err)

	evtIn, err = mesh.NewEvent("test", plEvtA, plEvtB, plEvtC)
	verify.NoError(t, err)
	data, err = json.Marshal(evtIn)
	verify.NoError(t, err)

	evtOut, err = mesh.NewEvent("empty")
	verify.NoError(t, err)
	err = json.Unmarshal(data, &evtOut)
	verify.NoError(t, err)
	verify.DeepEqual(t, evtOut, evtIn)
	pl := []*mesh.Event{}
	err = evtOut.Payload(&pl)
	verify.NoError(t, err)
	verify.DeepEqual(t, pl[0], plEvtA)
	verify.DeepEqual(t, pl[1], plEvtB)
	verify.DeepEqual(t, pl[2], plEvtC)
}
