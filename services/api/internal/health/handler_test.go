package health

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealth(t *testing.T) {
	for _, method := range []string{"GET", "POST", "DELETE"} {
		t.Run(method, func(t *testing.T) {
			w := httptest.NewRecorder()
			Handler(w, httptest.NewRequest(method, "/healthz", nil))
			if method != "GET" {
				if w.Code != 405 {
					t.Fatalf("status %d", w.Code)
				}
				return
			}
			var body struct {
				Status string `json:"status"`
				Models bool   `json:"electoral_models"`
			}
			if err := json.Unmarshal(w.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if w.Code != http.StatusOK || body.Status != "ok" || body.Models {
				t.Fatalf("unexpected response: %s", w.Body)
			}
			if w.Header().Get("Content-Type") != "application/json" {
				t.Fatal("missing content type")
			}
		})
	}
}
