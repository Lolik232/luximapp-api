package domain

import (
	"net/http"

	"github.com/Lolik232/luximapp-api/pkg/errors"
)

type HttpError struct {
	Code    uint   `json:"code,omitempty"`
	Message string `json:"err_msg,omitempty"`
}

var (
	codes = map[errors.ErrorType]uint{
		errors.NoType:             http.StatusInternalServerError,
		errors.ErrInvalidArgument: 1,
		errors.ErrDuplicateEntry:  http.StatusConflict,
	}
)

func New(err error) (*HttpError, int) {
	if err == nil {
		return &HttpError{
			Code:    codes[errors.NoType],
			Message: "Internal server error.",
		}, http.StatusInternalServerError
	}
	errtype := errors.GetType(err)
	msg := ""
	httpCode := 400

	switch errtype {
	//Users should not be aware of internal problems
	case errors.NoType:
		msg = "Internal server error,"
		httpCode = http.StatusInternalServerError
	case errors.ErrDuplicateEntry:
		msg = err.Error()
		httpCode = http.StatusConflict
	default:
		msg = err.Error()
	}

	return &HttpError{
		Code:    codes[errtype],
		Message: msg,
	}, httpCode
}
