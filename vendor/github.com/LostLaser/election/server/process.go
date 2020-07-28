package server

import (
	"time"
)

func (s *Server) run() {
	for {
		if s.state == running {
			if !s.pingMaster() || s.triggerElection {
				startElection(s)
				s.triggerElection = false
			}
		}
		time.Sleep(s.heartbeatPause)
	}
}

func (s *Server) pingMaster() bool {
	s.emitter.Write(s.id, s.master, "HEARTBEAT")
	if s.master == "" || (s.master != s.id && !s.NeighborServers[s.master].isUp()) {
		return false
	}

	return true
}

func (s *Server) setMaster(masterID string) {
	if !s.isUp() {
		return
	}
	s.electionLock.Lock()
	defer s.electionLock.Unlock()
	if masterID != s.id && s.id == s.master {
		s.emitter.Write(s.id, "", "NOT_MASTER")
	}
	s.emitter.Write(masterID, s.id, "ELECT")
	s.master = masterID
}

func (s *Server) isUp() bool {
	return s.state == running
}
