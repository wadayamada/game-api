package server

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"20dojo-online/pkg/db"
	"20dojo-online/pkg/http/middleware"
	"20dojo-online/pkg/server/handler"
	"20dojo-online/pkg/server/model"
)

// Serve HTTPサーバを起動する
func Serve(addr string) {

	rand.Seed(time.Now().UnixNano())

	m := model.NewModel(db.Conn)
	h := handler.NewHandler(m)

	/* ===== URLマッピングを行う ===== */
	http.HandleFunc("/setting/get", get(handler.HandleSettingGet()))
	http.HandleFunc("/user/create", post(handler.HandleUserCreate()))

	// middlewareは 20dojo-online/pkg/http/middleware パッケージを利用する
	http.HandleFunc("/user/get",
		get(middleware.Authenticate(handler.HandleUserGet())))
	http.HandleFunc("/user/update",
		post(middleware.Authenticate(handler.HandleUserUpdate())))

	http.HandleFunc("/gacha/draw",
		post(middleware.Authenticate(h.HandleGachaDraw())))

	http.HandleFunc("/ranking/list",
		get(middleware.Authenticate(handler.HandleRankingList())))

	http.HandleFunc("/collection/list",
		get(middleware.Authenticate(h.HandleCollectionList())))

	http.HandleFunc("/game/finish",
		post(middleware.Authenticate(handler.HandleGameFinish())))

	/* ===== サーバの起動 ===== */
	log.Println("Server running...")
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatalf("Listen and serve failed. %+v", err)
	}
}

// get GETリクエストを処理する
func get(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodGet)
}

// post POSTリクエストを処理する
func post(apiFunc http.HandlerFunc) http.HandlerFunc {
	return httpMethod(apiFunc, http.MethodPost)
}

// httpMethod 指定したHTTPメソッドでAPIの処理を実行する
func httpMethod(apiFunc http.HandlerFunc, method string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		// CORS対応
		writer.Header().Add("Access-Control-Allow-Origin", "*")
		writer.Header().Add("Access-Control-Allow-Headers", "Content-Type,Accept,Origin,x-token")

		// プリフライトリクエストは処理を通さない
		if request.Method == http.MethodOptions {
			return
		}
		// 指定のHTTPメソッドでない場合はエラー
		if request.Method != method {
			writer.WriteHeader(http.StatusMethodNotAllowed)
			writer.Write([]byte("Method Not Allowed"))
			return
		}

		// 共通のレスポンスヘッダを設定
		writer.Header().Add("Content-Type", "application/json")
		apiFunc(writer, request)
	}
}
