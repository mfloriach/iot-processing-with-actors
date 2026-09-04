package config

import "time"

const (
	NUM_CPUS       = 1
	NUM_OF_WORKERS = 2

	SENSOR_OF_ANALYSIS_ID = "0"
	MAILBOX_SIZE          = 10
	// QUANTUM               = 400

	PRINT_TELEMETRY   = 10 * time.Second
	SENSOR_PER_WORKER = 1_500_000
)
