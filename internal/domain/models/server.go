package models

type EncryptedCreateServerModel struct {
	UserID               string `json:"-"`
	EncryptedIp          string `json:"encrypted_ip"`
	EncryptedPort        string `json:"encrypted_port"`
	EncryptedDisplayName string `json:"encrypted_display_name"`
	IsMonitoringEnabled  bool   `json:"is_monitoring_enabled"`
	IsServerEnabled      bool   `json:"is_server_enabled"`
	CreatedAt            string `json:"created_at"`
}

type EncryptedServerModel struct {
	ID                   string `json:"id"`
	UserID               string `json:"-"`
	EncryptedIp          string `json:"encrypted_ip"`
	EncryptedPort        string `json:"encrypted_port"`
	EncryptedDisplayName string `json:"encrypted_display_name"`
	IsMonitoringEnabled  bool   `json:"is_monitoring_enabled"`
	IsServerEnabled      bool   `json:"is_server_enabled"`
	CreatedAt            string `json:"created_at"`
}
