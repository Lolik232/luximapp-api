package errors

import (
	"github.com/Lolik232/luximapp-api/internal/context_keys"
)

type LogError struct {
	ErrorMsg string                 `json:"error_msg"`
	Input    context_keys.InputData `json:"input_data"`
}
