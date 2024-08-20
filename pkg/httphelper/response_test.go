package httphelper

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
)

type FailingEncoder struct{}

var _ json.Marshaler = FailingEncoder{}

func (f FailingEncoder) MarshalJSON() ([]byte, error) {
	return nil, errors.New("failed to marshal json")
}

func TestWriteJSONResponse(t *testing.T) {

	testCases := []struct {
		desc     string
		status   int
		response interface{}
		expected string
	}{
		{
			desc:   "error response",
			status: 500,
			response: ErrorResponse{
				Error:  "internal server error",
				Status: 500,
			},
			expected: "{\"error\":\"internal server error\",\"status\":500}\n",
		},
		{
			desc:     "json error 500",
			status:   500,
			response: FailingEncoder{},
			expected: "",
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			h := httptest.NewRecorder()
			WriteJSONResponse(h, tC.status, tC.response)
			if h.Code != tC.status {
				t.Errorf("expected status code %d, got %d", tC.status, h.Code)
			}

			if h.Body.String() != tC.expected {
				t.Errorf("expected body %s, got %s", tC.expected, h.Body.String())
			}

			if h.Header().Get("Content-Type") != "application/json" {
				t.Errorf("expected Content-Type to be application/json, got %s", h.Header().Get("Content-Type"))
			}
		})
	}
}
