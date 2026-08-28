package main

import (
	"bytes"
	"datacollector/device"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"
)

func eventsHandler(manager *device.DeviceManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "SSE not supported", http.StatusInternalServerError)
			return
		}

		device := manager.GetDevice("sensor-1")
		if device == nil {
			http.Error(w, "device not found", http.StatusNotFound)
			return
		}

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var frame bytes.Buffer
		encoder := json.NewEncoder(&frame)

		for {
			select {
			case <-r.Context().Done():
				log.Println("client disconnected")
				return

			case <-ticker.C:
				frame.Reset()
				frame.Grow(128)
				frame.WriteString("data: ")

				if err := encoder.Encode(device.State()); err != nil {
					slog.Error("Error marshaling to JSON", slog.Any("error", err))
					continue
				}

				frame.WriteByte('\n')

				if _, err := w.Write(frame.Bytes()); err != nil {
					slog.Error("Error writing SSE frame", slog.Any("error", err))
					return
				}

				flusher.Flush()
			}
		}
	}
}

func StartServer(manager device.DeviceManager) {
	http.HandleFunc("/events", eventsHandler(&manager))

	log.Println("server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
