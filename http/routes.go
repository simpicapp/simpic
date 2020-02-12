package http

func (s *server) routes() {
	s.router.HandleFunc("/photo/{uuid}", s.handleGetPhoto()).Name("get_photo")
	s.router.HandleFunc("/photo", s.handleStorePhoto()).Methods("POST")
	s.router.HandleFunc("/timeline", s.handleTimeline())
}
