package model

import "database/sql"

type Model struct {
	db *sql.DB
}

// NewModel model構造体の初期化
func NewModel(db *sql.DB) *Model {
	return &Model{
		db: db,
	}
}
