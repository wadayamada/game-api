package model

import (
	"database/sql"
	"log"

	"20dojo-online/pkg/db"
)

type UserCollectionItem struct {
	UserID       string
	CollectionID string
}

// SelectUserCollectionsByPrimaryKey user_collection_itemテーブルをもとにユーザのアイテム所持情報を取得
func (m *Model) SelectUserCollectionsByUserID(id string) ([]*UserCollectionItem, error) {
	rows, err := db.Conn.Query("SELECT user_id,collection_item_id FROM user_collection_item where user_id=?", id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var userCollectionItems []*UserCollectionItem
	for rows.Next() {
		u := &UserCollectionItem{}
		if err := rows.Scan(&u.UserID, &u.CollectionID); err != nil {
			log.Println(err)
			return nil, err
		}
		userCollectionItems = append(userCollectionItems, u)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return userCollectionItems, nil
}

// InsertUserCollectionItem データベースをレコードを登録する
func InsertUserCollectionItem(user_id string, collection_item_id string) error {
	stmt, err := db.Conn.Prepare("INSERT INTO `user_collection_item`(`user_id`, `collection_item_id`) VALUES (?,?);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user_id, collection_item_id)
	return nil
}

// InsertUserCollectionItemTx トランザクションでuser_collection_itemテーブルにデータinsert
func InsertUserCollectionItemTx(user_id string, collection_item_id string, tx *sql.Tx) error {
	stmt, err := tx.Prepare("INSERT INTO `user_collection_item`(`user_id`, `collection_item_id`) VALUES (?,?);")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(user_id, collection_item_id)
	return nil
}

// InsertUserCollectionItemsTx トランザクションでuser_collection_itemテーブルにデータをまとめてbulk insert
func InsertUserCollectionItemsTx(uCISlice []*UserCollectionItem, tx *sql.Tx) error {
	queryString := "INSERT INTO `user_collection_item`(`user_id`, `collection_item_id`) VALUES"
	queryArgs := make([]interface{}, 0, len(uCISlice)*2)
	for i := range uCISlice {
		queryString += "(?,?),"
		queryArgs = append(queryArgs, uCISlice[i].UserID, uCISlice[i].CollectionID)
	}
	queryString = queryString[:len(queryString)-1]
	queryString += ";"
	stmtTx, err := tx.Prepare(queryString)
	if err != nil {
		return err
	}
	_, err = stmtTx.Exec(queryArgs...)
	return nil
}
