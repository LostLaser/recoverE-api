package server

// Election is an interface for all supported election types
type Election interface {
	start()
}
