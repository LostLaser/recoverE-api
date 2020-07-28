package server

type ringElection struct {
}

func (r ringElection) start(s *Server) {
	if isHighest(s) {
		notifyLow(s)
		s.setMaster(s.id)
	}
}
