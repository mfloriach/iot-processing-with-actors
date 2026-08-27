package main

import (
	"datacollector/device"
	"encoding/json"
	"fmt"
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

		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				log.Println("client disconnected")
				return

			case _ = <-ticker.C:
				state := manager.GetDevice("sensor-1").State()
				jsonData, err := json.Marshal(state)
				if err != nil {
					slog.Error("Error marshaling to JSON", slog.Any("error", err))
				}

				fmt.Fprintf(w, "data: %s\n\n", string(jsonData))

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
