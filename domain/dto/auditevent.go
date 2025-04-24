package dto

import "time"

type AuditEvent struct {
	Method     string    `json:"method"`
	Path       string    `json:"path"`
	RemoteAddr string    `json:"remote_addr"`
	Timestamp  time.Time `json:"timestamp"`
}
