//go:generate mockgen -source=$GOFILE -package=mock_$GOPACKAGE -destination=mock_$GOPACKAGE/$GOFILE
package model

import "database/sql"

type ModelInterface interface {
	SelectCollectionItemAll() ([]*CollectionItem, error)
	SelectUserCollectionsByUserID(string) ([]*UserCollectionItem, error)
}

type Model struct {
	db *sql.DB
}

//インターフェース満たしているかチェック
var _ ModelInterface = (*Model)(nil)

// NewModel model構造体の初期化
func NewModel(db *sql.DB) *Model {
	return &Model{
		db: db,
	}
}
