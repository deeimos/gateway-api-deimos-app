package models

type ServerMetric struct {
	CPUUsage      float64 `json:"cpu_usage"`
	MemoryUsage   float64 `json:"memory_usage"`
	DiskUsage     float64 `json:"disk_usage"`
	LoadAvg1      float64 `json:"load_avg_1"`
	LoadAvg5      float64 `json:"load_avg_5"`
	LoadAvg15     float64 `json:"load_avg_15"`
	NetworkRx     float64 `json:"network_rx"`
	NetworkTx     float64 `json:"network_tx"`
	DiskRead      float64 `json:"disk_read"`
	DiskWrite     float64 `json:"disk_write"`
	ProcessCount  int32   `json:"process_count"`
	IOWait        float64 `json:"io_wait"`
	UptimeSeconds int64   `json:"uptime_seconds"`
	Temperature   float64 `json:"temperature"`
	Status        string  `json:"status"`
	Timestamp     string  `json:"timestamp"`
}
