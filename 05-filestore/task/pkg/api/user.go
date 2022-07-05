package api

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func (s *server) getCurrentUserId(ctx context.Context, _ http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodGet {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	userId, err := getUserIdFrom(ctx)
	if err != nil {
		return nil, internalErr("get user id failed", err)
	}
	if userId == uuid.Nil {
		return nil, newApiErr(http.StatusUnauthorized, "unauthorized")
	}

	return userId, nil
}
