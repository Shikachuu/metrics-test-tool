package handler

import (
	"log/slog"
	"net/http"
	"testing"
)

func TestNewServer(t *testing.T) {
	testCases := []struct {
		desc   string
		logger *slog.Logger
		mux    *http.ServeMux
	}{
		{
			desc:   "success",
			logger: slog.Default(),
			mux:    http.NewServeMux(),
		},
		{
			desc:   "nil logger",
			logger: nil,
			mux:    http.NewServeMux(),
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			s := NewServer(tC.logger, tC.mux)

			if tC.logger == nil && s.logger != nil {
				t.Errorf("expected logger to be %v, got %v", tC.logger, s.logger)
			}

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected already registered panic, got nil")
				}
			}()

			s.mux.Handle("/health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

		})
	}
}

func TestServer_addRoutes(t *testing.T) {
	testCases := []struct {
		desc  string
		mux   *http.ServeMux
		route string
	}{
		{
			desc:  "success",
			mux:   http.NewServeMux(),
			route: "/health",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			t.Parallel()

			s := &Server{logger: nil, mux: tC.mux}

			s.addRoutes()

			defer func() {
				if r := recover(); r == nil {
					t.Errorf("expected already registered panic, got nil for route %s", tC.route)
				}
			}()

			tC.mux.Handle(tC.route, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))
		})
	}
}
