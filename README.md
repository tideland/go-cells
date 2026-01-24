# Tideland Go Cells

[![GitHub release](https://img.shields.io/github/release/tideland/go-cells.svg)](https://github.com/tideland/go-cells)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-cells/main/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-cells)](https://github.com/tideland/go-cells/blob/main/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/cells?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/cells?tab=packages)
![Workflow](https://github.com/tideland/go-cells/actions/workflows/build.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-cells)](https://goreportcard.com/report/tideland.dev/go/cells)

## Description

**Tideland Go Cells** provides a light-weight event-processing. It is realized
as a mesh of cells which can subsribe to each other. One cell can subscribe to
multiple cells as well as multiple cells can subscribe to one cell. Each cell
runs an individual developed and/or configured behavior with an own state.

I hope you like it. ;)

## Features

### Type-Safe Events (v0.2.0+)

Go Cells now supports type-safe event handling using Go generics:

```go
// Create type-safe events with metadata
type UserLogin struct {
    UserID   string
    Email    string
    LoginAt  time.Time
}

evt, err := mesh.NewEventTyped("user-login", UserLogin{
    UserID:  "user-123",
    Email:   "user@example.com",
    LoginAt: time.Now(),
},
    mesh.WithTraceID("trace-abc-123"),
    mesh.WithCorrelation("session-xyz-789"),
    mesh.WithCustom("ip", "192.168.1.1"),
)

// Extract typed payloads safely
user, err := mesh.PayloadAs[UserLogin](evt)
```

### Event Metadata

All events automatically include metadata:
- **ID**: Unique event identifier (auto-generated)
- **TraceID**: For distributed tracing
- **Correlation**: For correlating related events
- **Custom**: Application-specific key-value pairs

```go
meta := evt.Metadata()
fmt.Printf("Event ID: %s, Trace: %s\n", meta.ID, meta.TraceID)
```

### Backward Compatibility

The original `NewEvent()` and `Payload()` methods are still supported for gradual migration.

## Behaviors

The project already contains some standard behaviors, the number is still growing.

- **Aggregator** aggregates events and emits each aggregated value.
- **Broadcaster** simply emits received events to all subscribers.
- **Callback** calls a number of passed functions for each received event.
- **Collector** collects events which can be processed on demand.
- **Combo** waits for a user-defined combination of events.
- **Condition** tests events for conditions using a tester function and calls a
  processor then.
- **Countdown** counts a number of events down to zero and executes an event returning
  function. The event will be emitted then.
- **Counter** counts events, the counters can be retrieved.
- **Evaluator** evaluates events based on a user-defined function which returns a rating.
- **Filter** re-emits received events based on a user-defined filter. Those can be including
  or excluding.
- **Mapper** allows to analyse events and map them into new one for emitting.
- **One-Time** processes a user defined function only once for the first event, it will never
  called again. Outgoing events can be emitted during processing.
- **Pairer** allows to define a criterion for a first and second evend and a timeout
  between those. Matches and timouts will be emitted.
- **Rate Evaluator** measures times between a number of criterion fitting events and
  emits statistical data about these fittings.
- **Rate Window Evaluator** checks if a number of events in a given timespan matches
  a given criterion. In case it processes them.

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

