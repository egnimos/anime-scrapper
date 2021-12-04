package server_engine

var (
	ParsedServers *ServerList
	Servers Server = &server{}
)

type Server interface {
	Server(parsed ParseServer)
}

type server struct{}

func (s *server) Server(parsed ParseServer) {
	ParsedServers = parsed.ParseServers()
}
