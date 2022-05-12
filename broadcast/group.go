package broadcast

type Group struct {
	stop        chan struct{}
	publish     chan interface{}
	subscribe   chan chan interface{}
	unsubscribe chan chan interface{}
	subscribers map[chan interface{}]struct{}
}

func NewGroup() *Group {
	return &Group{
		make(chan struct{}, 0),
		make(chan interface{}, 1),
		make(chan chan interface{}, 1),
		make(chan chan interface{}, 1),
		make(map[chan interface{}]struct{}),
	}
}

func (group *Group) Start() {
	for {
		select {
		case <-group.stop:
			return
		case msg := <-group.publish:
			for channel := range group.subscribers {
				// Try sending, skip if blocking
				select {
				case channel <- msg:
				default:
				}
			}
		case channel := <-group.subscribe:
			group.subscribers[channel] = struct{}{}
		case channel := <-group.unsubscribe:
			delete(group.subscribers, channel)
		}
	}
}

func (group *Group) Stop() {
	close(group.stop)
}

func (group *Group) Subscribe(channel chan interface{}) {
	group.subscribe <- channel
}

func (group *Group) Unsubscribe(channel chan interface{}) {
	group.unsubscribe <- channel
}

func (group *Group) Publish(message interface{}) {
	group.publish <- message
}

func (group *Group) Len() int {
	return len(group.subscribe)
}
