package main

type ChatRoom struct {
	name  string
	users []*ChatUser
	calls chan syncRoomCall
}

func NewChatRoom(name string) *ChatRoom {
	return &ChatRoom{
		name,
		make([]*ChatUser, 0),
		make(chan syncRoomCall),
	}
}

func (room *ChatRoom) HandleCalls() {
	for {
		call := <-room.calls
		if _, ok := call.(stop); !ok {
			call.execute(room)
		} else {
			return
		}
	}
}

type syncRoomCall interface {
	execute(*ChatRoom)
}

func (room *ChatRoom) Stop() {
	room.calls <- stop{}
}

type stop struct{}

func (call stop) execute(_ *ChatRoom) {}

func (room *ChatRoom) Add(user *ChatUser) {
	room.calls <- addUser{user}
}

type addUser struct {
	user *ChatUser
}

func (call addUser) execute(room *ChatRoom) {
	room.users = append(room.users, call.user)
}

func (room *ChatRoom) Remove(user *ChatUser) {
	room.calls <- removeUser{user}
}

type removeUser struct {
	user *ChatUser
}

func (call removeUser) execute(room *ChatRoom) {
	for i, u := range room.users {
		if u == call.user {
			room.users = append(room.users[:i], room.users[i+1:]...)
		}
	}
}

func (room *ChatRoom) Send(msg string) {
	room.calls <- sendMessage{msg}
}

type sendMessage struct {
	msg string
}

func (call sendMessage) execute(room *ChatRoom) {
	for _, user := range room.users {
		user.Send(call.msg)
	}
}
