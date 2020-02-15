package http

import "net/http"

func (s *server) routes() {
	s.router.HandleFunc("/login", s.handleAuthenticate()).Methods("POST")
	s.router.HandleFunc("/thumbnail/{uuid}", s.handleGetThumbnail())
	s.router.HandleFunc("/photo/{uuid}", s.handleGetPhoto()).Name("get_photo")
	s.router.HandleFunc("/photo", s.handleStorePhoto()).Methods("POST")
	s.router.HandleFunc("/timeline", s.handleTimeline())

	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(s.staticDir)))
}
