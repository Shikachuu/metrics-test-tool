package handler

import (
	"net/http"

	"github.com/Shikachuu/metrics-test-tool/pkg/httphelper"
)

func (s *Server) buildHealthHandler() http.HandlerFunc {
	type HealthResponse struct {
		Status string `json:"status"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		httphelper.WriteJSONResponse(w, http.StatusOK, HealthResponse{Status: "ok"})
	}
}
