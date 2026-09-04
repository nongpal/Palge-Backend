package data

import (
	"database/sql"
	"time"
)

type AuditLog struct {
	ID         int64     `json:"id"`
	UserID     *int64    `json:"user_id"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID *int64    `json:"resource_id"`
	IPAddress  string    `json:"ip_address"`
	UserAgent  string    `json:"user_agent"`
	Result     string    `json:"result"`
	CreatedAt  time.Time `json:"created_at"`
}

type AuditLogModel struct {
	DB *sql.DB
}
