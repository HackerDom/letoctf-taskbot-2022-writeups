package api

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (s *server) auth(registerUser bool) apiHandler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
		if r.Method != http.MethodPost {
			return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
		}

		userId, err := getUserIdFrom(ctx)
		if err != nil {
			return nil, internalErr("get user id from session failed", err)
		}
		if userId != uuid.Nil {
			return nil, nil
		}

		var creds struct {
			Username string
			Pass     string
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&creds); err != nil {
			return nil, errorf(http.StatusBadRequest, "", "decode json failed: %v", err, exposeErrToResp)
		}

		pass := s.getMAC([]byte(creds.Pass))
		if registerUser {
			userId, err = s.dbStore.CreateUser(ctx, creds.Username, pass, uuid.Nil)
			if err != nil {
				return nil, internalErr("create user failed", err)
			}
		} else {
			userId, err = s.dbStore.GetUserId(ctx, creds.Username, pass)
			if err != nil {
				return nil, newApiErr(http.StatusUnauthorized, "bad credentials")
			}
		}

		session, err := s.sessionStore.Get(r, userSessionKey)
		if err != nil {
			return nil, internalErr("get session failed", err)
		}

		session.Values[userIdSessionValueKey] = userId.String()
		if err := session.Save(r, w); err != nil {
			return nil, internalErr("save session failed", err)
		}

		return nil, nil
	}
}

func (s *server) getMAC(message []byte) []byte {
	mac := hmac.New(sha256.New, s.hmacKey)
	mac.Write(message)
	return mac.Sum(nil)
}

func (s *server) logout(_ context.Context, w http.ResponseWriter, r *http.Request) (interface{}, *errApi) {
	if r.Method != http.MethodPost {
		return nil, newApiErr(http.StatusMethodNotAllowed, "method not allowed")
	}

	session, err := s.sessionStore.Get(r, userSessionKey)
	if err != nil {
		return nil, internalErr("get session failed", err)
	}
	if session.IsNew {
		return nil, newApiErr(http.StatusUnauthorized, "unauthorized")
	}

	session.Options.MaxAge = -1
	if err := session.Save(r, w); err != nil {
		return nil, internalErr("save session failed", err)
	}

	return nil, nil
}
