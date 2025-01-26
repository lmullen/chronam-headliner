package headliner

// Routes registers the handlers for the URLs that should be served.
func (a *App) Routes() {
	a.Router.HandleFunc("/chronamurl", a.ChronamUrlHandler()).Methods("POST")
}
