package broadcast

type Server struct {
	stop     chan struct{}
	groups   map[string]*Group
	requests chan groupRequest
}
type groupRequest struct {
	name  string
	reply chan *Group
}

func NewServer() *Server {
	return &Server{
		make(chan struct{}, 0),
		make(map[string]*Group),
		make(chan groupRequest, 1),
	}
}

func (server *Server) Start() {
	for {
		select {
		case <-server.stop:
			return
		case req := <-server.requests:
			group, exists := server.groups[req.name]
			if !exists {
				group = NewGroup()
				go group.Start()
				server.groups[req.name] = group
			}
			req.reply <- group
		}
	}
}

func (server *Server) Stop() {
	close(server.stop)
}

func (server *Server) GetGroup(name string) *Group {
	reply := make(chan *Group)
	server.requests <- groupRequest{name, reply}
	return <-reply
}
