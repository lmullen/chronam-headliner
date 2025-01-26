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

// The App type shares access to resources.
type App struct {
	Server   *http.Server
	Config   Config
	Router   *mux.Router
	AIClient *AIClient
}

func NewApp(ctx context.Context) *App {
	a := App{}

	a.Config.Address = "0.0.0.0" + ":" + "8050"

	// Create the router, store it in the struct, initialize the routes, and
	// register the middleware.
	router := mux.NewRouter()
	a.Router = router
	a.Routes()
	// s.Middleware()

	a.Server = &http.Server{
		Addr:         a.Config.Address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	a.AIClient = NewAIClient(ctx)

	return &a
}

func (a *App) Run() error {

	slog.Info("starting the server", "address", "http://"+a.Config.Address)
	err := a.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown closes the connection to the database and shutsdown the server.
func (a *App) Shutdown() {
	slog.Info("shutting down the web server")
	err := a.Server.Shutdown(context.TODO())
	if err != nil {
		slog.Error("error shutting down web server", "error", err)
	}
}
