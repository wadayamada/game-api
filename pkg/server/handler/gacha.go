package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/db"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type gachaDrawRequest struct {
	Times int32 `json:"times"`
}

type gachaDrawResult struct {
	CollectionID string `json:"collectionID"`
	Name         string `json:"name"`
	Rarity       int32  `json:"rarity"`
	IsNew        bool   `json:"isNew"`
}

type gachaDrawResponse struct {
	Results []gachaDrawResult `json:"results"`
}

// HandleGachaDraw リクエストされた回数分、ガチャを引く処理
func (h *Handler) HandleGachaDraw() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyからガチャを引く回数を取得
		var requestBody gachaDrawRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		// timesが1以上かバリデーション
		if requestBody.Times < 1 {
			log.Println(errors.New("times is less than 1"))
			response.BadRequest(writer, fmt.Sprintf("times is less than 1. times=%v", requestBody.Times))
			return
		}

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// トランザクション開始
		tx, err := db.Conn.Begin()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		user, err := model.SelectUserByUserIDTx(userID, tx)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			tx.Rollback()
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			tx.Rollback()
			return
		}
		if constant.GachaCoinConsumption*requestBody.Times > user.Coin {
			log.Println(errors.New("not enough coins"))
			response.BadRequest(writer, fmt.Sprintf("not enough coins"))
			tx.Rollback()
			return
		}

		// collection_itemテーブルから全てのアイテム情報を取得
		collectionItems, err := h.model.SelectCollectionItemAll()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			tx.Rollback()
			return
		}

		// user_collection_itemテーブルをもとにユーザのアイテム所持情報を取得
		userCollectionItems, err := model.SelectUserCollectionsByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			tx.Rollback()
			return
		}

		// gacha_probabilityテーブルの全件取得
		gachaProbabilities, err := model.SelectGachaProbabilityAll()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			tx.Rollback()
			return
		}

		sumRatio := 0
		for _, v := range gachaProbabilities {
			sumRatio += v.Ratio
		}

		userCollectionItemMap := make(map[string]struct{}, len(userCollectionItems)+int(requestBody.Times))
		for _, v := range userCollectionItems {
			userCollectionItemMap[v.CollectionID] = struct{}{}
		}
		collectionItemMap := make(map[string]*model.CollectionItem, len(collectionItems))
		for i, v := range collectionItems {
			collectionItemMap[v.CollectionID] = collectionItems[i]
		}

		gachaDrawResults := make([]gachaDrawResult, 0, requestBody.Times)
		newUserCollectionItems := make([]*model.UserCollectionItem, 0, int(requestBody.Times))
		for i := 0; i < int(requestBody.Times); i++ {
			// 乱数取って1キャラ選択する
			choiceNumber := rand.Intn(sumRatio)
			tmpRate := 0
			var selectedCollectionID string
			for _, v := range gachaProbabilities {
				tmpRate += v.Ratio
				if tmpRate > choiceNumber {
					selectedCollectionID = v.CollectionItemID
					break
				}
			}

			_, hasItem := userCollectionItemMap[selectedCollectionID]
			selectedCollectionItem := collectionItemMap[selectedCollectionID]
			gachaDrawResults = append(gachaDrawResults, gachaDrawResult{
				CollectionID: selectedCollectionItem.CollectionID,
				Name:         selectedCollectionItem.Name,
				Rarity:       selectedCollectionItem.Rarity,
				IsNew:        !hasItem,
			})

			// user_collection_itemテーブルにinsertするためにIsNewのデータをスライスに集める
			if !hasItem {
				newUserCollectionItems = append(newUserCollectionItems, &model.UserCollectionItem{
					UserID:       userID,
					CollectionID: selectedCollectionID,
				})
				userCollectionItemMap[selectedCollectionID] = struct{}{}
			}
		}

		user.Coin -= constant.GachaCoinConsumption * requestBody.Times
		if err = model.UpdateUserByUserIDTx(user, tx); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			tx.Rollback()
			return
		}
		// まとめてuser_collection_itemテーブルにinsertする
		if len(newUserCollectionItems) > 0 {
			if err := model.InsertUserCollectionItemsTx(newUserCollectionItems, tx); err != nil {
				log.Println(err)
				response.InternalServerError(writer, "Internal Server Error")
				tx.Rollback()
				return
			}
		}

		//トランザクション終了
		if err := tx.Commit(); err != nil {
			log.Printf("Commit is failed :%v", err)
		}

		// レスポンスに必要な情報を詰めて返却
		response.Success(writer, &gachaDrawResponse{Results: gachaDrawResults})
	}
}
