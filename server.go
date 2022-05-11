package main

// Global instance
var chatServer = ChatServer{
	make(map[string]*ChatRoom),
	make(chan syncCall),
}

type ChatServer struct {
	rooms map[string]*ChatRoom
	calls chan syncCall
}

// A call to a server's method which requires synchronization
type syncCall interface {
	execute(*ChatServer)
}

func (server *ChatServer) HandleRequests() {
	for {
		call := <-server.calls
		call.execute(server)
	}
}

func (server *ChatServer) GetChatRoom(name string) *ChatRoom {
	call := getRoom{
		name,
		make(chan *ChatRoom),
	}
	server.calls <- call
	return <-call.response
}

type getRoom struct {
	name     string
	response chan *ChatRoom
}

func (call getRoom) execute(server *ChatServer) {
	room, ok := server.rooms[call.name]
	if !ok {
		room = &ChatRoom{
			call.name,
			make([]*ChatUser, 0),
		}
		server.rooms[call.name] = room
	}
	call.response <- room
}

func (server *ChatServer) DeleteChatRoom(name string) {
	server.calls <- deleteRoom{
		name,
	}
}

type deleteRoom struct {
	name string
}

func (call deleteRoom) execute(server *ChatServer) {
	room, ok := server.rooms[call.name]
	if ok && len(room.users) == 0 {
		delete(server.rooms, call.name)
	}
}
