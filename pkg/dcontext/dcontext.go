package dcontext

import (
	"context"
)

type key string

const (
	userIDKey key = "userID"
)

// SetUserID ContextへユーザIDを保存する
func SetUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

// GetUserIDFromContext ContextからユーザIDを取得する
func GetUserIDFromContext(ctx context.Context) string {
	var userID string
	if ctx.Value(userIDKey) != nil {
		userID = ctx.Value(userIDKey).(string)
	}
	return userID
}
