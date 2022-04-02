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

const (
	MethodGET = "GET"
)

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

func Get[T any](url string, opts ...func(*requestConfig)) (*T, error) {
	config := initConfig()

	ctx := config.ctx
	defer config.cancelFunc()

	req, err := http.NewRequestWithContext(ctx, MethodGET, url, nil)
	if err != nil {
		return nil, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}
	defer res.Body.Close()

	var out T
	err = json.NewDecoder(res.Body).Decode(&out)
	if err != nil {
		return nil, fmt.Errorf("req: GET request to %s failed: %w", url, err)
	}

	return &out, nil
}

func WithTimeout(t time.Duration) func(r *requestConfig) {
	return func(r *requestConfig) {
		r.ctx, r.cancelFunc = context.WithTimeout(r.ctx, t)
	}
}

func WithContext(c context.Context) func(r *requestConfig) {
	return func(r *requestConfig) {
		r.ctx, r.cancelFunc = context.WithCancel(c)
	}
}
