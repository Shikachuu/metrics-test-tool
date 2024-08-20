package handler

import (
	"log/slog"
	"net/http"
)

// Server is the main handler for the application and implements the http.Handler interface
// to handle all incoming requests. It has all the http.HandlerFuncs registered to the mux.
type Server struct {
	logger *slog.Logger
	mux    *http.ServeMux
}

// Server implements the http.Handler interface and handles all incoming requests.
var _ http.Handler = &Server{}

// NewServer creates a new server instance with the given logger.
func NewServer(logger *slog.Logger, mux *http.ServeMux) *Server {
	s := &Server{logger: logger, mux: mux}
	s.addRoutes()

	return s
}

// addRoutes registers all the routes to the mux.
func (s *Server) addRoutes() {
	s.mux.HandleFunc("/health", s.buildHealthHandler())
}

// ServeHTTP handles all incoming requests and pass them down to the underlying mux.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
