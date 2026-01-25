// Copyright 2010-2026 Tideland / Frank Mueller. All rights reserved.
// Use of this source code is governed by the BSD 3-Clause
// license that can be found in the LICENSE file.

package mesh

// BehaviorFunc simplifies implementation of a behavior when only
// one function is needed. It can be deployed via
//
//	myMesh.Go("my-name", BehaviorFunc(myFunc))
type BehaviorFunc func(cell Cell, in Receptor, out Emitter) error

// Go implements Behavior.
func (bf BehaviorFunc) Go(cell Cell, in Receptor, out Emitter) error {
	return bf(cell, in, out)
}

// RequestFunc defines a function signature for the request
// behavior. It is called per received event.
type RequestFunc func(cell Cell, evt *Event, out Emitter) error

// RequestBehavior is a simple behavior using a function
// to process the received events.
type RequestBehavior struct {
	rf RequestFunc
}

// NewRequestBehavior creates a behavior based on the given
// processing function.
func NewRequestBehavior(rf RequestFunc) RequestBehavior {
	return RequestBehavior{
		rf: rf,
	}
}

// Go implements Behavior.
func (rb RequestBehavior) Go(cell Cell, in Receptor, out Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			if err := rb.rf(cell, evt, out); err != nil {
				return err
			}
		}
	}
}
