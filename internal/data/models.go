package data

import (
	"database/sql"
)

type Models struct {
	MessageModel MessageModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		MessageModel: MessageModel{DB: db},
	}
}
