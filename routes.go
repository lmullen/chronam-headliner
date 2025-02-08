package headliner

// Routes registers the handlers for the URLs that should be served.
func (a *App) Routes() {
	a.Router.HandleFunc("/chronamurl", a.ChronamUrlHandler()).Methods("POST")
	a.Router.HandleFunc("/", a.RootHandler()).Methods("GET")

	// Log 404 errors
	a.Router.NotFoundHandler = loggingMiddleware(a.NotFoundHandler())
}
