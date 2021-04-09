package handler

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"20dojo-online/pkg/dcontext"
	"20dojo-online/pkg/server/model"
	"20dojo-online/pkg/server/model/mock_model"

	"github.com/golang/mock/gomock"
)

func TestHandler_HandleCollectionList(t *testing.T) {
	type args struct {
		userID string
	}
	tests := []struct {
		name     string
		mock     func(*gomock.Controller) *mock_model.MockModelInterface
		args     args
		wantCode int
		wantBody interface{}
	}{
		{
			name: "正常系_ユーザがコレクションを2つ持っている",
			mock: func(ctrl *gomock.Controller) *mock_model.MockModelInterface {
				mockModel := mock_model.NewMockModelInterface(ctrl)
				mockModel.
					EXPECT().
					SelectCollectionItemAll().
					Return([]*model.CollectionItem{
						{CollectionID: "sample_collectionID1", Name: "sample_collection1", Rarity: 1},
						{CollectionID: "sample_collectionID2", Name: "sample_collection2", Rarity: 2},
					},
						nil)
				mockModel.
					EXPECT().
					SelectUserCollectionsByUserID("sample_userID").
					Return([]*model.UserCollectionItem{
						{UserID: "sample_userID", CollectionID: "sample_collectionID1"},
						{UserID: "sample_userID", CollectionID: "sample_collectionID2"},
					},
						nil)
				return mockModel
			},
			args:     args{"sample_userID"},
			wantCode: 200,
			wantBody: `{"collections":[{"collectionID":"sample_collectionID1","name":"sample_collection1","rarity":1,"hasItem":true},{"collectionID":"sample_collectionID2","name":"sample_collection2","rarity":2,"hasItem":true}]}`,
		},
		{
			name: "異常系_selectCollectonでエラー",
			mock: func(ctrl *gomock.Controller) *mock_model.MockModelInterface {
				mockModel := mock_model.NewMockModelInterface(ctrl)
				mockModel.
					EXPECT().
					SelectCollectionItemAll().
					Return(nil,
						errors.New("Internal Server Error"))
				return mockModel
			},
			args:     args{"sample_userID"},
			wantCode: 500,
			wantBody: `{"code":500,"message":"Internal Server Error"}`,
		},
		{
			name: "異常系_selectUserCollectionでエラー",
			mock: func(ctrl *gomock.Controller) *mock_model.MockModelInterface {
				mockModel := mock_model.NewMockModelInterface(ctrl)
				mockModel.
					EXPECT().
					SelectCollectionItemAll().
					Return([]*model.CollectionItem{
						{CollectionID: "sample_collectionID1", Name: "sample_collection1", Rarity: 1},
						{CollectionID: "sample_collectionID2", Name: "sample_collection2", Rarity: 2},
					},
						nil)
				mockModel.
					EXPECT().
					SelectUserCollectionsByUserID("sample_userID").
					Return(nil,
						errors.New("Internal Server Error"))
				return mockModel
			},
			args:     args{"sample_userID"},
			wantCode: 500,
			wantBody: `{"code":500,"message":"Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			modelMock := tt.mock(ctrl)
			h := NewHandler(modelMock)
			req := httptest.NewRequest(http.MethodPost, "/collection/list", nil)
			req.Header.Add("accept", "application/json")
			req.Header.Add("x-token", "sample_token")
			// レスポンスを受け止める*httptest.ResponseRecorder
			got := httptest.NewRecorder()
			//実際にはmiddleware.Authenticateで行うユーザID取ってくる処理のモック
			ctx := req.Context()
			ctx = dcontext.SetUserID(ctx, tt.args.userID)
			h.HandleCollectionList()(got, req.WithContext(ctx))
			body, err := ioutil.ReadAll(got.Result().Body)
			if err != nil {
				t.Errorf("ioutil.ReadAll failed")
			}
			if got.Result().StatusCode != tt.wantCode {
				t.Errorf("status code = %d, want %d", got.Result().StatusCode, tt.wantCode)
			}
			if string(body) != tt.wantBody {
				t.Errorf("response body = %v\n, want %v\n", string(body), tt.wantBody)
			}
		})
	}
}
