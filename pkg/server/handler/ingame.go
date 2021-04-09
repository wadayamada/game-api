package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"20dojo-online/pkg/constant"
	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type gameFinishRequest struct {
	Score int `json:"score"`
}
type gameFinishResponse struct {
	Coin int `json:"coin"`
}

// HandleGameFinish インゲーム終了API
func HandleGameFinish() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyからスコア情報を取得
		var requestBody gameFinishRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}
		//scoreが0以上かバリデーション
		if requestBody.Score < 0 {
			log.Println(errors.New("score is minus"))
			response.BadRequest(writer, fmt.Sprintf("score is minus. score=%v", requestBody.Score))
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
		//ユーザIDからデータを取得
		user, err := model.SelectUserByUserID(userID)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}
		if user == nil {
			log.Println(errors.New("user not found"))
			response.BadRequest(writer, fmt.Sprintf("user not found. userID=%s", userID))
			return
		}

		//コインの追加処理
		coin := requestBody.Score * constant.RewardCoinRate
		user.Coin += int32(coin)
		if user.HighScore < int32(requestBody.Score) {
			user.HighScore = int32(requestBody.Score)
		}
		//userテーブルでコインとハイスコア更新
		if err = model.UpdateUserCoinHighScoreByUserID(user); err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// 獲得したコイン情報を返却
		response.Success(writer, &gameFinishResponse{
			Coin: coin,
		})
	}
}
