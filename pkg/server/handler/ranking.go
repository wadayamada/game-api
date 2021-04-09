package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

type ranking struct {
	UserID   string `json:"userId"`
	UserName string `json:"userName"`
	Rank     int32  `json:"rank"`
	Score    int32  `json:"score"`
}

type rankingListResponse struct {
	Ranks []ranking `json:"ranks"`
}

// HandleRankingList ランキング情報取得API
func HandleRankingList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		//ランキング表示のスタートランクを取り出す
		startRank, err := strconv.Atoi(request.URL.Query().Get("start"))
		if err != nil {
			log.Println(err)
			response.BadRequest(writer, fmt.Sprintf("get request query not found"))
			return
		}
		//startRankが1以上かバリデーション
		if startRank < 1 {
			log.Println(errors.New("starkRank < 1"))
			response.BadRequest(writer, fmt.Sprintf("starkRank < 1. startRank=%v", startRank))
			return
		}

		//startRankから始まるランキングデータをuserテーブルから取得
		userSlice, err := model.SelectUsersByRanking(startRank)
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		rankingSlice := make([]ranking, 0, len(userSlice))
		sR := startRank
		for _, user := range userSlice {
			r := ranking{
				UserID:   user.ID,
				UserName: user.Name,
				Rank:     int32(sR),
				Score:    user.HighScore,
			}
			rankingSlice = append(rankingSlice, r)
			sR++
		}
		// レスポンスに必要な情報を詰めて返却
		response.Success(writer, &rankingListResponse{Ranks: rankingSlice})
	}
}
