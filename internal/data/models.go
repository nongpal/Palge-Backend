package data

import (
	"database/sql"
)

type Models struct {
	Accounts    AccountModel
	Permissions PermissionModel
	Tokens      TokenModel
	Users       UserModel
	AuditLogs   AuditLogModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Accounts:    AccountModel{DB: db},
		Permissions: PermissionModel{DB: db},
		Tokens:      TokenModel{DB: db},
		Users:       UserModel{DB: db},
		AuditLogs:   AuditLogModel{DB: db},
	}
}
