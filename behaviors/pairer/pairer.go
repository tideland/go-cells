// Tideland Go Cells - Behaviors - Pairer
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package pairer // import "tideland.dev/go/cells/behaviors/pairer"

//--------------------
// IMPORTS
//--------------------

import (
	"time"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// CONSTANTS
//--------------------

const (
	TopicMatch   = "pairer:match"
	TopicTimeout = "pairer:timeout"
)

//--------------------
// HELPER
//--------------------

// PairerFunc is used by the pair behavior and has to return true, if
// the given event matches a criterion. The first event is the already
// found first hit. So it could be nil and the criterion test has to be
// only for the second.
type PairerFunc func(fstEvt, sndEvt *mesh.Event) bool

// Pair contains the paired events or a possible timeout.
type Pair struct {
	First  *mesh.Event
	Second *mesh.Event
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior provides a behavior checking a stream of incoming event for
// matching a wanted pair crterion. This is defined as pairer function
// when starting the behavior. A duration defines the maximum allowed
// duration between first and second event.
type Behavior struct {
	matches  PairerFunc
	duration time.Duration
	hit      *mesh.Event
}

var _ mesh.Behavior = (*Behavior)(nil)

// New creates a new instance of the pairer behavior.
func New(pairer PairerFunc, duration time.Duration) *Behavior {
	return &Behavior{
		matches:  pairer,
		duration: duration,
	}
}

// Go implements the mesh.Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	var ticker *time.Ticker
	quitec := make(<-chan time.Time)
	tickc := quitec
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if !b.matches(b.hit, evt) {
				continue
			}
			if b.hit == nil {
				// First hit resets ticker.
				b.hit = evt
				ticker = time.NewTicker(b.duration)
				tickc = ticker.C
				continue
			}
			// Second hit.
			pair := Pair{
				First:  b.hit,
				Second: evt,
			}
			b.hit = nil
			tickc = quitec
			ticker.Stop()
			out.Emit(TopicMatch, pair)
		case <-tickc:
			// Timeout!
			out.Emit(TopicTimeout)
			ticker.Stop()
		}
	}
}

// EOF
