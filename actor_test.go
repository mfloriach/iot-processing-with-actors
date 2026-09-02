package main

import (
	"datacollector/device"
	"testing"
	"time"
)

func TestActorStatePublishesAtomicSnapshot(t *testing.T) {
	actor := NewActor("sensor-1", device.Dispatch)

	first := device.Telemetry{
		Sample: device.Sample{
			DeviceID: "sensor-1",
			TTL:      time.Now().Add(-time.Second),
		},
		Temperature: 21.5,
		Humidity:    40.1,
	}

	if !actor.Send(first) {
		t.Fatal("expected first send to enqueue the actor")
	}

	if !actor.Update(1) {
		t.Fatal("expected first update to consume the queued message")
	}

	if actor.Update(1) {
		t.Fatal("expected empty mailbox to stop the actor")
	}

	snapshot := actor.State()
	if snapshot.Data.Temperature != first.Temperature {
		t.Fatalf("unexpected snapshot temperature: got %v want %v", snapshot.Data.Temperature, first.Temperature)
	}

	second := device.Telemetry{
		Sample: device.Sample{
			DeviceID: "sensor-1",
			TTL:      time.Now().Add(-time.Second),
		},
		Temperature: 22.75,
		Humidity:    39.4,
	}

	if !actor.Send(second) {
		t.Fatal("expected second send to enqueue the actor")
	}

	if !actor.Update(1) {
		t.Fatal("expected second update to consume the queued message")
	}

	if actor.Update(1) {
		t.Fatal("expected empty mailbox to stop the actor")
	}

	updated := actor.State()
	if updated.Data.Temperature != second.Temperature {
		t.Fatalf("unexpected updated temperature: got %v want %v", updated.Data.Temperature, second.Temperature)
	}

	if snapshot.Data.Temperature != first.Temperature {
		t.Fatalf("previous snapshot was mutated: got %v want %v", snapshot.Data.Temperature, first.Temperature)
	}
}
