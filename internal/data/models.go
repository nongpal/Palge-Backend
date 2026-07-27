package data

import "database/sql"

type Models struct {
	Accounts AccountModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Accounts: AccountModel{DB: db},
	}
}
