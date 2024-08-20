package handler

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestServer_buildHealthHandler(t *testing.T) {
	testCases := []struct {
		desc     string
		status   int
		response string
	}{
		{
			desc:     "health handler",
			status:   200,
			response: "{\"status\":\"ok\"}\n",
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			res := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/health", nil)
			s := &Server{mux: http.NewServeMux()}
			s.addRoutes()

			s.ServeHTTP(res, req)

			if res.Code != tC.status {
				t.Errorf("expected status code %d, got %d", tC.status, res.Code)
			}

			if res.Body.String() != tC.response {
				t.Errorf("expected body %s, got %s", tC.response, res.Body.String())
			}
		})
	}
}
