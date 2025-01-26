package headliner

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	Address string // The address at which this will be hosted, e.g.: localhost:8090
}

// The Server type shares access to resources.
type Server struct {
	Server *http.Server
	Config Config
	Router *mux.Router
}

func NewServer(ctx context.Context) *Server {
	s := Server{}

	// Read the configuration from environment variables. The `getEnv()` function
	// will provide a default.
	s.Config.Address = "0.0.0.0" + ":" + "8050"

	// Create the router, store it in the struct, initialize the routes, and
	// register the middleware.
	router := mux.NewRouter()
	s.Router = router
	s.Routes()
	// s.Middleware()

	s.Server = &http.Server{
		Addr:         s.Config.Address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      s.Router,
	}

	return &s
}

func (s *Server) Run() error {
	slog.Info("starting the server", "address", "http://"+s.Config.Address)
	err := s.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown closes the connection to the database and shutsdown the server.
func (s *Server) Shutdown() {
	slog.Info("shutting down the web server")
	err := s.Server.Shutdown(context.TODO())
	if err != nil {
		slog.Error("error shutting down web server", "error", err)
	}
}
