package handlers

import (
	"backend/middleware"
	"context"
)

func getUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	return userID, ok
}
