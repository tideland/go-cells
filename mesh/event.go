// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

package mesh // import "tideland.dev/go/cells/mesh"

//--------------------
// IMPORTS
//--------------------

import (
	"tideland.dev/go/cells/mesh/internal"
)

//--------------------
// RE-EXPORTS
//--------------------

// Event represents an event with topic, payload, and metadata.
type Event = internal.Event

// Metadata contains event metadata for tracing and correlation.
type Metadata = internal.Metadata

// MetadataOption is a functional option for configuring event metadata.
type MetadataOption = internal.MetadataOption

// Event creation functions.
var (
	NewEvent             = internal.NewEvent
	NewEventWithMetadata = internal.NewEventWithMetadata
)

// NewEventTyped creates a type-safe Event with a typed payload and metadata options.
// This is the preferred method for creating events with type safety.
func NewEventTyped[T any](topic string, payload T, opts ...MetadataOption) (*Event, error) {
	return internal.NewEventTyped(topic, payload, opts...)
}

// PayloadAs returns the typed payload of the event with compile-time type safety.
// This is the preferred method for accessing typed payloads.
func PayloadAs[T any](evt *Event) (T, error) {
	return internal.PayloadAs[T](evt)
}

// Metadata option functions.
var (
	WithID          = internal.WithID
	WithTraceID     = internal.WithTraceID
	WithCorrelation = internal.WithCorrelation
	WithCustom      = internal.WithCustom
)

// EOF
