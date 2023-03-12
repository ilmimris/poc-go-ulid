package constant

import (
	"context"
	"time"
)

type ContextKey string

const (
	// Context Key
	ContextKeyLatency   ContextKey = "latency"
	ContextKeyRequestID ContextKey = "request_id"
	ContextKeyErrorResp ContextKey = "error_response"
)

func (c ContextKey) ToString() string {
	return string(c)
}

func (c ContextKey) Error(ctx context.Context) error {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(error)
	}
	return nil
}

func (c ContextKey) String(ctx context.Context) (string, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(string), nil
	}
	return "", ErrCtxDataNotExist
}

func (c ContextKey) Int(ctx context.Context) (int, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(int), nil
	}
	return 0, ErrCtxDataNotExist
}

func (c ContextKey) Int64(ctx context.Context) (int64, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(int64), nil
	}
	return 0, ErrCtxDataNotExist
}

func (c ContextKey) Bool(ctx context.Context) (bool, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(bool), nil
	}
	return false, ErrCtxDataNotExist
}

func (c ContextKey) Float64(ctx context.Context) (float64, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(float64), nil
	}
	return 0, ErrCtxDataNotExist
}

func (c ContextKey) Time(ctx context.Context) (time.Time, error) {
	if ctx.Value(c) != nil {
		return ctx.Value(c).(time.Time), nil
	}
	return time.Time{}, ErrCtxDataNotExist
}
