package models

import "time"

type ServerModel struct {
	ID                  string
	UserID              string
	IP                  string
	Port                string
	DisplayName         string
	IsMonitoringEnabled bool
	CreatedAt           time.Time
}

type EncryptedServerModel struct {
	ID                   string
	UserID               string
	EncryptedIP          string
	EncryptedPort        string
	EncryptedDisplayName string
	IsMonitoringEnabled  bool
	CreatedAt            time.Time
}

type CreateServerModel struct {
	UserID              string
	IP                  string
	Port                string
	DisplayName         string
	IsMonitoringEnabled bool
}

type EncryptedCreateServerModel struct {
	UserID               string
	EncryptedIP          string
	EncryptedPort        string
	EncryptedDisplayName string
	IsMonitoringEnabled  bool
}

type EncryptedServerInfo struct {
	EncryptedIP          string
	EncryptedPort        string
	EncryptedDisplayName string
}

type DecryptedServerInfo struct {
	IP          string
	Port        string
	DisplayName string
}
