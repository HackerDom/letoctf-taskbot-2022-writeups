package api

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
)

type errApi struct {
	error

	httpCode int
	msg      string
}

type option func(*errApi)

func errorf(httpCode int, msg, format string, err error, opts ...option) *errApi {
	newErr := &errApi{
		error:    fmt.Errorf(format, err),
		httpCode: httpCode,
		msg:      msg,
	}
	for _, opt := range opts {
		opt(newErr)
	}

	return newErr
}

func newApiErr(code int, msg string) *errApi {
	return errorf(code, "", "%v", errors.New(msg), exposeErrToResp)
}

func internalErr(msg string, err error) *errApi {
	return errorf(http.StatusInternalServerError, msg, fmt.Sprintf("%v: %v", msg, "%v"), err)
}

func exposeErrToResp(err *errApi) {
	err.msg = err.Error()
}

func middlewareErr(w http.ResponseWriter, logger *zerolog.Logger, err error, msg string) {
	w.WriteHeader(http.StatusInternalServerError)
	logger.Err(err).Msg(msg)

	if _, printErr := fmt.Fprintf(w, "{\"response\":\"%v\"}", msg); printErr != nil {
		logger.Err(printErr).Err(err).Msg("write response failed")
	}
}
