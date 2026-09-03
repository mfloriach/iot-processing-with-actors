# iot_data_collection

A Go telemetry simulator for exploring actor mailboxes, sharded workers, backpressure, and SSE streaming under load.

It creates a large set of synthetic sensors, feeds them random telemetry, and keeps the latest state in memory for fast reads. The project is useful as a concurrency playground or as a benchmark for experimenting with scheduling and queueing strategies.

## What It Does

- Simulates `NUM_OF_SENSORS` devices and continuously generates telemetry for them.
- Uses a per-device mailbox so each device is the authority over its own state.
- Requeues busy devices through a shared ready queue so one hot device does not monopolize a worker.
- Publishes the latest state of `sensor-0` over Server-Sent Events.
- Tracks latency, backlog, GC, and backpressure metrics while the app runs.

## Requirements

- Go `1.26.4` or newer.
- A terminal that can run long-lived Go processes.
- Optional: `curl` or a browser if you want to inspect the SSE stream.

## Quick Start

From a clone of this repository:

```bash
$ go run .
```

The process starts the simulator, opens:

- `http://localhost:8080/events` for SSE updates
- `http://localhost:6060/debug/pprof/` for pprof

To watch the event stream:

```bash
$ curl -N -H "Accept: text/event-stream" http://localhost:8080/events
```

Expected output is a stream of `data:` frames containing JSON snapshots, for example:

```text
data: {"Data":{"DeviceID":"sensor-0","Temperature":0,"Humidity":0,"Battery":0,"Noise":0},"Online":true}
```

## Configuration

Runtime settings live in [`config/configuration.go`](config/configuration.go):

- `NUM_CPUS` - sets `GOMAXPROCS`
- `NUM_OF_WORKERS_UPDATING` - number of workers that drain ready actors
- `NUM_OF_WORKERS_GENERETIC_NOISE` - number of workers that generate synthetic telemetry
- `NUM_OF_SENSORS` - total number of simulated sensors
- `SENSOR_OF_ANALYSIS_ID` - the device exposed through `/events`

If you want to change the scale of the simulation, update those constants and rerun the program.

## How It Works

The main flow is:

1. `main.go` initializes metrics, creates the device manager, and registers the sensors.
2. `scheduler.Scheduler` generates telemetry and sends it into the device manager.
3. Each `device.Actor` owns a mailbox and publishes a lock-free snapshot for readers.
4. `server.go` serves the latest snapshot for `sensor-0` over SSE.
5. `mesures/latency_states.go` prints periodic latency and memory statistics.

The current design is intentionally in-memory and ephemeral. It does not persist telemetry to disk or a database.

## Project Status

This repository is still evolving. The current roadmap in `TODO.md` includes:

- alarm priority over telemetry
- work stealing
- per-device ordering guarantees
- coalescing duplicate updates
- rate limiting per device
- deadline-aware scheduling
