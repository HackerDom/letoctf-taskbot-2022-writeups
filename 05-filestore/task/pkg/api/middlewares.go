package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/rs/zerolog"
)

type apiHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errApi)
type middleHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request)
type middleware func(middleHandler) middleHandler

func corsMiddleware(next middleHandler) middleHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		zerolog.Ctx(ctx).Info().Msg("request")

		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.WriteHeader(http.StatusOK)
			return
		}

		next(ctx, w, r)
	}
}

func loggerMiddleware(logger *zerolog.Logger) middleware {
	return func(next middleHandler) middleHandler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			l := logger.With().
				Str("url", r.URL.String()).
				Str("method", r.Method).
				Logger()
			next(l.WithContext(ctx), w, r)
		}
	}
}

func authMiddleware(sessionStore sessions.Store) middleware {
	return func(next middleHandler) middleHandler {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
			logger := zerolog.Ctx(ctx)

			session, err := sessionStore.Get(r, userSessionKey)
			if err != nil {
				logger.Err(err).Msg("get session failed")
				next(ctx, w, r)
				return
			}

			if session.IsNew {
				next(ctx, w, r)
				return
			}

			rawUserId, ok := session.Values[userIdSessionValueKey].(string)
			if !ok {
				middlewareErr(w, logger, err, "get session value failed")
				return
			}
			userId, err := uuid.Parse(rawUserId)
			if err != nil {
				middlewareErr(w, logger, err, "parse user id failed")
				return
			}

			ctx = withUserId(ctx, userId)
			next(ctx, w, r)
		}
	}
}

type userIdContextKey struct{}

func withUserId(ctx context.Context, userId uuid.UUID) context.Context {
	return context.WithValue(ctx, userIdContextKey{}, userId)
}

func getUserIdFrom(ctx context.Context) (uuid.UUID, error) {
	v := ctx.Value(userIdContextKey{})
	if v == nil {
		return uuid.Nil, nil
	}

	userId, ok := v.(uuid.UUID)
	if !ok {
		panic(any("userId must have type uuid.UUID"))
	}

	return userId, nil
}

func responseMiddleware(handler apiHandler) middleHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(ctx)

		type Resp struct {
			Response interface{} `json:"response"`
		}
		var resp Resp

		apiResp, apiErr := handler(ctx, w, r)
		switch {
		case apiErr != nil:
			if apiErr.httpCode == http.StatusInternalServerError {
				logger.Err(apiErr).Msg("handle request failed")
			} else {
				logger.Info().Err(apiErr).Msg("got api error")
			}
			w.WriteHeader(apiErr.httpCode)
			resp = Resp{
				Response: apiErr.msg,
			}
		case apiResp == nil:
			resp = Resp{
				Response: "successful",
			}
		default:
			resp = Resp{
				Response: apiResp,
			}
		}

		jsonResp, err := json.Marshal(resp)
		if err != nil {
			logger.Err(err).Msg("marshal response to json failed")
		}

		var respData []byte

		if data, ok := apiResp.([]byte); ok && apiErr == nil {
			respData = data
		} else {
			respData = jsonResp
		}

		n, err := w.Write(respData)
		if err != nil {
			logger.Err(err).Msg("write response failed")
		}
		if n != len(respData) {
			logger.Err(err).
				Str("want write", fmt.Sprint(len(respData))).
				Str("written", fmt.Sprint(n)).
				Msg("response wasn't write completely")
		}

		logger.Info().Msg("response")
	}
}

func collectMiddlewares(middlewares ...middleware) middleware {
	return func(handler middleHandler) middleHandler {
		m := handler
		for i := len(middlewares) - 1; i >= 0; i-- {
			m = middlewares[i](m)
		}

		return m
	}
}

func toHttpHandler(ctx context.Context, handler middleHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(ctx, 10 * time.Second)
		defer cancel()
		handler(ctx, w, r)
	}
}
