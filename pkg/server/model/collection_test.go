package model

import (
	"database/sql"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSQLConnect_SelectCollectionItemAll(t *testing.T) {
	tests := []struct {
		name    string
		dbMock  func() (*sql.DB, sqlmock.Sqlmock, error)
		want    []*CollectionItem
		wantErr bool
	}{
		{
			"正常系_コレクションデータ2つ取得",
			func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Error("error occured when opening a mock database connection", err)
					return db, mock, err
				}
				mock.ExpectQuery("SELECT id, name, rarity FROM collection_item").
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "rarity"}).
						AddRow("sample_collectionID1", "sample_collection1", 1).
						AddRow("sample_collectionID2", "sample_collection2", 2))
				return db, mock, err
			},
			[]*CollectionItem{
				{CollectionID: "sample_collectionID1", Name: "sample_collection1", Rarity: 1},
				{CollectionID: "sample_collectionID2", Name: "sample_collection2", Rarity: 2},
			},
			false,
		},
		{
			"異常系_selectErr",
			func() (*sql.DB, sqlmock.Sqlmock, error) {
				db, mock, err := sqlmock.New()
				if err != nil {
					t.Error("error occured when opening a mock database connection", err)
					return db, mock, err
				}
				mock.ExpectQuery("SELECT id, name, rarity FROM collection_item").
					WillReturnError(fmt.Errorf("test database error"))
				return db, mock, err
			},
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, _, err := tt.dbMock()
			if err != nil {
				t.Error("error occured when opening a mock database connection", err)
				return
			}
			u := &Model{
				db: db,
			}
			defer db.Close()
			got, err := u.SelectCollectionItemAll()
			if (err != nil) != tt.wantErr {
				t.Errorf("SQLConnect.SelectCollectionItemAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SQLConnect.SelectCollectionItemAll() = %v, want %v", got, tt.want)
				return
			}
		})
	}
}
