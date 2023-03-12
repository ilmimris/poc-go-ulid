package constant

import "errors"

var (
	// Error Context
	ErrCtxDataNotExist = errors.New("the data not exist")

	ErrNotSupportedDisplayPayload = errors.New("not supported to display payload")
)
