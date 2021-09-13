# Tideland Go Cells

[![GitHub release](https://img.shields.io/github/release/tideland/go-cells.svg)](https://github.com/tideland/go-cells)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-cells/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-cells)](https://github.com/tideland/go-cells/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/cells?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/cells?tab=packages)
![Workflow](https://github.com/tideland/go-cells/actions/workflows/go.yml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-cells)](https://goreportcard.com/report/tideland.dev/go/cells)

## Description

**Tideland Go Cells** provides a light-weight event-processing. It is realized
as a mesh of cells which can subsribe to each other. One cell can subscribe to
multiple cells as well as multiple cells can subscribe to one cell. Each cell
runs an individual developed and/or configured behavior with an own state.

I hope you like it. ;)

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

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

