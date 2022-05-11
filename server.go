package main

// Global instance
var chatServer = ChatServer{
	make(map[string]*ChatRoom),
	make(chan syncServerCall),
}

type ChatServer struct {
	rooms map[string]*ChatRoom
	calls chan syncServerCall
}

// A call to a server's method which requires synchronization
type syncServerCall interface {
	execute(*ChatServer)
}

func (server *ChatServer) HandleCalls() {
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
		room = NewChatRoom(call.name)
		go room.HandleCalls()
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
		room.Stop()
		delete(server.rooms, call.name)
	}
}
