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
		serverResponseCode int
	}{
		{
			name: "simple GET request", requestMethod: req.MethodGET, serverResponseCode: 200,
		},
	}

	for _, c := range cases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(c.serverDelay)
			json.NewEncoder(w).Encode(&responseType{Success: true})
		}))

		if c.requestMethod == req.MethodGET {
			res, err := req.Get[responseType](server.URL)
			if err != nil {
				t.Errorf("did not expect error, but got one: %s", err.Error())
			}
			if !res.Res().Success {
				t.Errorf("wanted success = true, got %v", res.Res())
			}

			if res.StatusCode() != c.serverResponseCode {
				t.Errorf("wanted response code %d, got %d", c.serverResponseCode, res.StatusCode())
			}
		}
	}
}
