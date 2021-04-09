package model

import (
	"log"
)

type CollectionItem struct {
	CollectionID string
	Name         string
	Rarity       int32
}

// SelectCollectionItemAll collection_itemテーブルから全てのアイテム情報を取得
func (u *Model) SelectCollectionItemAll() ([]*CollectionItem, error) {
	var collectionItems []*CollectionItem
	rows, err := u.db.Query("SELECT id, name, rarity FROM collection_item")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		c := &CollectionItem{}
		if err := rows.Scan(&c.CollectionID, &c.Name, &c.Rarity); err != nil {
			log.Println(err)
			return nil, err
		}
		collectionItems = append(collectionItems, c)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return collectionItems, nil
}
