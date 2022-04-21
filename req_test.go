package req_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/preslavrachev/req"
)

type responseType struct {
	Success bool `json:"success,omitempty"`
}

func TestReq(t *testing.T) {

	cases := []struct {
		name               string
		requestMethod      string
		serverDelay        time.Duration
		opts               []req.RequestOption
		serverResponseCode int
		expectedErr        bool
	}{
		{
			name:               "simple GET request",
			requestMethod:      http.MethodGet,
			serverResponseCode: 200,
		},
		{
			name:          "simple GET request with timeout",
			requestMethod: http.MethodGet,
			serverDelay:   time.Minute,
			opts:          []req.RequestOption{req.WithTimeout(time.Nanosecond)},
			expectedErr:   true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				time.Sleep(c.serverDelay)
				json.NewEncoder(w).Encode(&responseType{Success: true})
			}))

			if c.requestMethod == http.MethodGet {
				res, err := req.Get[responseType](server.URL, c.opts...)

				switch e, _ := err, c.expectedErr; {
				case e != nil, true:
					return
				case e != nil, false:
					t.Errorf("did not expect error, but got one: %s", err.Error())
				case e == nil, true:
					t.Fatal("expected err but did not get one")
				}

				if !res.Res().Success {
					t.Errorf("wanted success = true, got %v", res.Res())
				}

				if res.StatusCode() != c.serverResponseCode {
					t.Errorf("wanted response code %d, got %d", c.serverResponseCode, res.StatusCode())
				}
			}
		})
	}
}
