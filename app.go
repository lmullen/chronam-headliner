package headliner

import (
	"context"
	_ "embed"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/gorilla/mux"
)

// The Config type stores configuration which is read from environment variables.
type Config struct {
	Address string // The address at which this will be hosted, e.g.: localhost:8090
}

// The App type shares access to resources.
type App struct {
	Server      *http.Server
	Config      Config
	Router      *mux.Router
	AIClient    *anthropic.Client
	ShutdownCtx context.Context
	MakePrompt  PromptMaker
	Store       *sync.Map
}

func NewApp(ctx context.Context) (*App, error) {
	a := App{}

	a.ShutdownCtx = ctx

	a.Config.Address = "0.0.0.0" + ":" + "8050"

	var m sync.Map
	a.Store = &m

	// Create the router, store it in the struct, initialize the routes, and
	// register the middleware.
	router := mux.NewRouter()
	a.Router = router
	a.Middleware()
	a.Routes()

	a.Server = &http.Server{
		Addr:         a.Config.Address,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      a.Router,
	}

	a.AIClient = anthropic.NewClient()

	p, err := MakePromptTemplate()
	if err != nil {
		slog.Error("error creating prompt template", "error", err)
		return nil, err
	}
	a.MakePrompt = p

	return &a, nil
}

func (a *App) Run() error {

	slog.Info("starting the web server", "address", "http://"+a.Config.Address)
	err := a.Server.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}
	return err
}

// Shutdown closes the connection to the database and shutsdown the server.
func (a *App) Shutdown() {
	slog.Info("shutting down the web server")
	err := a.Server.Shutdown(a.ShutdownCtx)
	if err != nil {
		slog.Error("error shutting down web server", "error", err)
	}
}
