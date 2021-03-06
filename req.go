package req

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type responseFormat int

const (
	ResponseFormatJSON responseFormat = iota
)

type RequestOption func(*requestConfig)

type requestConfig struct {
	ctx        context.Context
	cancelFunc context.CancelFunc
}

var noOpCancelFunc = func() { /* do nothing */ }

func initConfig() *requestConfig {
	return &requestConfig{
		ctx:        context.Background(),
		cancelFunc: noOpCancelFunc,
	}
}

type ResultTransformer[T any] interface {
	Do() T
}

type HTTPResult[T any] struct {
	r *T
	c int
}

func (r *HTTPResult[T]) Res() *T         { return r.r }
func (r *HTTPResult[T]) StatusCode() int { return r.c }

func Get[T any](url string, opts ...RequestOption) (HTTPResult[T], error) {
	return req[T](http.MethodGet, url, opts...)
}

func req[T any](method string, url string, opts ...RequestOption) (HTTPResult[T], error) {
	config := initConfig()
	r := HTTPResult[T]{}

	for _, opt := range opts {
		opt(config)
	}

	ctx := config.ctx
	defer config.cancelFunc()

	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return r, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return r, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}
	defer res.Body.Close()

	r.c = res.StatusCode

	var out T
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return r, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}

	r.r = &out
	return r, nil
}

func WithTimeout(t time.Duration) RequestOption {
	return func(r *requestConfig) {
		r.ctx, r.cancelFunc = context.WithTimeout(r.ctx, t)
	}
}

func WithContext(c context.Context) RequestOption {
	return func(r *requestConfig) {
		r.ctx, r.cancelFunc = context.WithCancel(c)
	}
}
