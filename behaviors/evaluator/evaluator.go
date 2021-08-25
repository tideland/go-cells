// Tideland Go Cells - Behaviors - Evaluator
//
// Copyright (C) 2010-2021 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package evaluator // import "tideland.dev/go/cells/behaviors/evaluator"

//--------------------
// IMPORTS
//--------------------

import (
	"sort"

	"tideland.dev/go/cells/mesh"
)

//--------------------
// TOPICS
//--------------------

const (
	TopicEvaluate       = "evaluate!"
	TopicEvaluationDone = "evaluation-done"
	TopicReset          = "reset!"
	TopicResetDone      = "reset-done"
)

//--------------------
// HELPER
//--------------------

// EvaluationFunc is a function returning a rating for each received event.
type EvaluationFunc func(evt *mesh.Event) (float64, error)

// Evaluation contains the aggregated result of all evaluations.
type Evaluation struct {
	Count     int
	MinRating float64
	MaxRating float64
	AvgRating float64
	MedRating float64
}

//--------------------
// BEHAVIOR
//--------------------

// Behavior ebvaluates incomming events nummerically.
type Behavior struct {
	evaluate      EvaluationFunc
	maxRatings    int
	ratings       []float64
	sortedRatings []float64
}

var _ mesh.Behavior = &Behavior{}

// New initializes and returns a new Behavior using the given function
// for the evaluation of the individual received events.
func New(evaluate EvaluationFunc) *Behavior {
	return &Behavior{
		evaluate: evaluate,
	}
}

// Go implements the Behavior interface.
func (b *Behavior) Go(cell mesh.Cell, in mesh.Receptor, out mesh.Emitter) error {
	for {
		select {
		case <-cell.Context().Done():
			return nil
		case evt := <-in.Pull():
			switch evt.Topic() {
			case TopicReset:
				b.maxRatings = 0
				b.ratings = nil
				b.sortedRatings = nil
				out.Emit(TopicResetDone)
			case TopicEvaluate:
				evaluation := b.evaluateRatings()
				out.Emit(TopicEvaluationDone, evaluation)
			default:
				rating, err := b.evaluate(evt)
				if err != nil {
					return err
				}
				b.ratings = append(b.ratings, rating)
				if b.maxRatings > 0 && len(b.ratings) > b.maxRatings {
					// Take care for size.
					b.ratings = b.ratings[1:]
				}
				if len(b.sortedRatings) < len(b.ratings) {
					// Let it grow up to the needed size.
					b.sortedRatings = append(b.sortedRatings, 0.0)
				}
			}
		}
	}
}

// evaluateRatings evaluates the collected ratings.
func (b *Behavior) evaluateRatings() Evaluation {
	var evaluation Evaluation

	copy(b.sortedRatings, b.ratings)
	sort.Float64s(b.sortedRatings)
	// Count.
	evaluation.Count = len(b.sortedRatings)
	// Average.
	totalRating := 0.0
	for _, rating := range b.sortedRatings {
		totalRating += rating
	}
	evaluation.AvgRating = totalRating / float64(evaluation.Count)
	// Median.
	if evaluation.Count%2 == 0 {
		// Even, have to calculate.
		middle := evaluation.Count / 2
		evaluation.MedRating = (b.sortedRatings[middle-1] + b.sortedRatings[middle]) / 2
	} else {
		// Odd, can take the middle.
		evaluation.MedRating = b.sortedRatings[evaluation.Count/2]
	}
	// Minimum and maximum.
	evaluation.MinRating = b.sortedRatings[0]
	evaluation.MaxRating = b.sortedRatings[len(b.sortedRatings)-1]

	return evaluation
}

// EOF
