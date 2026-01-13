package handlers

import (
	"backend/middleware"
	"context"
)

func getUserIDFromContext(ctx context.Context) (int64, bool) { // Retrieves the user ID from the context
	userID, ok := ctx.Value(middleware.UserIDKey).(int64)
	return userID, ok
}
