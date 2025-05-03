package models

type ServerForecast struct {
	Timestamp               string  `json:"timestamp"`
	CpuLoad                 float64 `json:"cpu_load"`
	MemoryLoad              float64 `json:"memory_load"`
	DiskUsage               float64 `json:"disk_usage"`
	LoadAvg1                float64 `json:"load_avg_1"`
	LoadAvg5                float64 `json:"load_avg_5"`
	NetworkRx               float64 `json:"network_rx"`
	NetworkTx               float64 `json:"network_tx"`
	AvailabilityProbability float64 `json:"availability_probability"`
	Status                  string  `json:"status"`
}
