package models

import "time"

type ServerMetric struct {
	CPUUsage      float64
	MemoryUsage   float64
	DiskUsage     float64
	LoadAvg1      float64
	LoadAvg5      float64
	LoadAvg15     float64
	NetworkRX     float64
	NetworkTX     float64
	DiskRead      float64
	DiskWrite     float64
	ProcessCount  int32
	IOWait        float64
	UptimeSeconds int64
	Temperature   float64
	Status        string
	Timestamp     time.Time
}

type CreateMetricModel struct {
	ServerID      string
	CPUUsage      float64
	MemoryUsage   float64
	DiskUsage     float64
	LoadAvg1      float64
	LoadAvg5      float64
	LoadAvg15     float64
	NetworkRX     float64
	NetworkTX     float64
	DiskRead      float64
	DiskWrite     float64
	ProcessCount  int32
	IOWait        float64
	UptimeSeconds int64
	Temperature   float64
	Status        string
}
