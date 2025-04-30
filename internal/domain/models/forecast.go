package models

import "time"

type ServerForecast struct {
	CpuLoad                 float64
	MemoryLoad              float64
	DiskUsage               float64
	LoadAvg_1               float64
	LoadAvg_5               float64
	NetworkRx               float64
	NetworkTx               float64
	AvailabilityProbability float64
	Timestamp               time.Time
	Status                  string
}
