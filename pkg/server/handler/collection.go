package handler

import (
	"errors"
	"log"
	"net/http"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// CollectionItem レスポンスの中身の構造体
type CollectionItem struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	HasItem      bool   `json:"hasItem"`
}

type collectionListResponse struct {
	Collections []CollectionItem `json:"collections"`
}

// HandleCollectionList ユーザの所持キャラクター情報を返すAPI
func (h *Handler) HandleCollectionList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		//collection_itemテーブルから全てのアイテム情報を取得
		itemInfos, err := h.model.SelectCollectionItemAll()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		//user_collection_itemテーブルをもとにユーザのアイテム所持情報を取得
		userCollectionItems, err := model.SelectUserCollectionsByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		userCollectionItemMap := make(map[string]struct{}, len(userCollectionItems))
		for i := range userCollectionItems {
			userCollectionItemMap[userCollectionItems[i].CollectionID] = struct{}{}
		}

		// レスポンスに必要な情報を詰めて返却
		collectionItemSlice := make([]CollectionItem, 0, len(itemInfos))
		for i := range itemInfos {

			c := CollectionItem{}

			_, c.HasItem = userCollectionItemMap[itemInfos[i].CollectionID]
			c.CollectionID = itemInfos[i].CollectionID
			c.Name = itemInfos[i].Name
			c.Rarity = itemInfos[i].Rarity

			collectionItemSlice = append(collectionItemSlice, c)
		}
		response.Success(writer, &collectionListResponse{Collections: collectionItemSlice})

	}
}
