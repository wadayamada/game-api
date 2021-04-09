package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/http/response"
	"20dojo-online/pkg/server/model"
)

// HandleUserCreate ユーザ情報作成処理
func HandleUserCreate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyから更新後情報を取得
		var requestBody userCreateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
			return
		}

		// UUIDでユーザIDを生成する
		userID, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// UUIDで認証トークンを生成する
		authToken, err := uuid.NewRandom()
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// データベースにユーザデータを登録する
		err = model.InsertUser(&model.User{
			ID:        userID.String(),
			AuthToken: authToken.String(),
			Name:      requestBody.Name,
			HighScore: 0,
			Coin:      0,
		})
		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		// 生成した認証トークンを返却
		response.Success(writer, &userCreateResponse{Token: authToken.String()})
	}
}

type userCreateRequest struct {
	Name string `json:"name"`
}

type userCreateResponse struct {
	Token string `json:"token"`
}

// HandleUserGet ユーザ情報取得処理
func HandleUserGet() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// Contextから認証済みのユーザIDを取得
		ctx := request.Context()
		userID := dcontext.GetUserIDFromContext(ctx)
		if userID == "" {
			log.Println(errors.New("userID is empty"))
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		var user *model.User
		var err error

		user, err = model.SelectUserByUserID(userID)

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

		// レスポンスに必要な情報を詰めて返却
		response.Success(writer, &userGetResponse{
			ID:        user.ID,
			Name:      user.Name,
			HighScore: user.HighScore,
			Coin:      user.Coin,
		})
	}
}

type userGetResponse struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	HighScore int32  `json:"highScore"`
	Coin      int32  `json:"coin"`
}

// HandleUserUpdate ユーザ情報更新処理
func HandleUserUpdate() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// リクエストBodyから更新後情報を取得
		var requestBody userUpdateRequest
		if err := json.NewDecoder(request.Body).Decode(&requestBody); err != nil {
			log.Println(err)
			response.BadRequest(writer, "Bad Request")
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

		var user *model.User
		var err error

		//ユーザIDからデータを取得
		user, err = model.SelectUserByUserID(userID)

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

		//名前をrequestの名前に書き換える
		user.Name = requestBody.Name
		err = model.UpdateUserByUserID(user)

		if err != nil {
			log.Println(err)
			response.InternalServerError(writer, "Internal Server Error")
			return
		}

		response.Success(writer, nil)
	}
}

type userUpdateRequest struct {
	Name string `json:"name"`
}
