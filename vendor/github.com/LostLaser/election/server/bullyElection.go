package server

type bullyElection struct {
}

func startElection(s *Server) {
	s.emitter.Write(s.id, "", "ELECTION_STARTED")
	if isHighest(s) {
		notifyLow(s)
		s.setMaster(s.id)
		s.emitter.Write(s.id, "", "ELECTED")
	}
	s.emitter.Write(s.id, "", "ELECTION_ENDED")
}

func isHighest(s *Server) bool {
	for id, neighbor := range s.NeighborServers {
		if id > s.id {
			if neighbor.isUp() {
				s.emitter.Write(s.id, id, "START_NEW_ELECTION")
				neighbor.triggerElection = true
				return false
			}
		}
	}
	return true
}

func notifyLow(s *Server) {
	for key, neighbor := range s.NeighborServers {
		if key < s.id {
			neighbor.setMaster(s.id)
		}
	}
}
