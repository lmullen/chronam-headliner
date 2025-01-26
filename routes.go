package headliner

// Routes registers the handlers for the URLs that should be served.
func (s *Server) Routes() {
	s.Router.HandleFunc("/chronamurl", s.ChronamUrlHandler()).Methods("POST")
}
