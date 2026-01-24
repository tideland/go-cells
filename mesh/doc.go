// Tideland Go Cells - Mesh
//
// Copyright (C) 2010-2026 Frank Mueller / Tideland / Oldenburg / Germany
//
// All rights reserved. Use of this source code is governed
// by the new BSD license.

// Package mesh is the runtime package of the Tideland cells event processing.
// It provides types for meshed cells running individual behaviors.
//
// These behaviors are defined based on an interface and can be added to the
// mesh. Here they are running concurrently and can be networked to communicate
// via events. Several useful behaviors are already provided with the behaviors
// package.
//
// # Creating a Mesh
//
// New meshes are created with
//
//     ctx := context.Background()
//     msh := mesh.New(ctx)
//
// and cells are started with
//
//    msh.Go("foo", NewFooBehavior())
//    msh.Go("bar", NewBarBehavior())
//    msh.Go("baz", NewBazBehavior())
//
// These cells can subscribe each other with
//
//    msh.Subscribe("foo", "bar")
//    msh.Subscribe("foo", "baz")
//
// so that events which are emitted by the cell "foo" will be
// received by the cells "bar" and "baz". Each cell can subscribe
// to multiple other subscribers and even circular subscriptions are
// no problem. But handle with care.
//
// # Events and Type Safety
//
// Events from the outside are emitted using
//
//     msh.Emit("foo", "topic", 42)
//
// For type-safe event creation with metadata, use NewEventTyped:
//
//     evt, err := mesh.NewEventTyped("user-login", userPayload,
//         mesh.WithTraceID("trace-123"),
//         mesh.WithCorrelation("session-456"),
//         mesh.WithCustom("source", "api-gateway"),
//     )
//
// For type-safe payload extraction, use PayloadAs:
//
//     user, err := mesh.PayloadAs[User](evt)
//
// # Metadata
//
// All events automatically include metadata with a unique ID. You can add:
// - TraceID: For distributed tracing
// - Correlation: For correlating related events
// - Custom fields: For application-specific metadata
//
// Access metadata with:
//
//     meta := evt.Metadata()
//     fmt.Println(meta.ID, meta.TraceID, meta.Custom)
//
// # Emitters
//
// In case of many emits to one cell you can get an emitter
// with
//
//     emtr, err := msh.Emitter("foo")
//
// and
//
//     emtr.Emit("answer", 42)
//
package mesh // import "tideland.dev/go/cells/mesh"

// EOF
