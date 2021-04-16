# Tideland Go Cells

[![GitHub release](https://img.shields.io/github/release/tideland/go-cells.svg)](https://github.com/tideland/go-cells)
[![GitHub license](https://img.shields.io/badge/license-New%20BSD-blue.svg)](https://raw.githubusercontent.com/tideland/go-cells/master/LICENSE)
[![Go Module](https://img.shields.io/github/go-mod/go-version/tideland/go-cells)](https://github.com/tideland/go-cells/blob/master/go.mod)
[![GoDoc](https://godoc.org/tideland.dev/go/together?status.svg)](https://pkg.go.dev/mod/tideland.dev/go/together?tab=packages)
[![Workflow](https://img.shields.io/github/workflow/status/tideland/go-cells/build)](https://github.com/tideland/go-cells/actions/)
[![Go Report Card](https://goreportcard.com/badge/github.com/tideland/go-cells)](https://goreportcard.com/report/tideland.dev/go/together)

## Description

**Tideland Go Cells** provides a light-weight event-processing. It is realized
as a mesh of cells which can subsribe to each other. One cell can subscribe to
multiple cells as well as multiple cells can subscribe to one cell. Each cell
runs an individual developed and/or configured behavior with an own state.

I hope you like it. ;)

## Behaviors

The project already contains some standard behaviors, the number is still growing.

* **Aggregator** aggregates events and emits each aggregated value.
* **Broadcaster** simply emits received events to all subscribers.
* **Callback** calls a number of passed functions for each received event.
* **Collector** collects events which can be processed on demand.
* **Combo** waits for a user-defined combination of events.
* **Condition** tests events for conditions using a tester function and calls a processor then.

## Contributors

- Frank Mueller (https://github.com/themue / https://github.com/tideland / https://tideland.dev)

