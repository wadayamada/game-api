package model

import (
	"log"

	"20dojo-online/pkg/db"
)

type GachaProbability struct {
	CollectionItemID string
	Ratio            int
}

func SelectGachaProbabilityAll() ([]*GachaProbability, error) {
	rows, err := db.Conn.Query("SELECT collection_item_id,ratio FROM gacha_probability")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()
	var gachaProbabilitySlice []*GachaProbability
	for rows.Next() {
		gachaProb := &GachaProbability{}
		if err := rows.Scan(&gachaProb.CollectionItemID, &gachaProb.Ratio); err != nil {
			log.Println(err)
			return nil, err
		}
		gachaProbabilitySlice = append(gachaProbabilitySlice, gachaProb)
	}
	if err := rows.Err(); err != nil {
		log.Println(err)
		return nil, err
	}
	return gachaProbabilitySlice, nil
}
