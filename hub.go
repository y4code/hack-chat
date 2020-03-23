package main

// Hub 维持了一些 开好的房间
type Hub struct {
	// 不同的聊天室 room
	rooms map[string]*Room
}

// maybe a hub maker?
func newHub() *Hub {
	return &Hub{
		rooms: make(map[string]*Room),
	}
}

//func (hub *Hub) run() {
//
//}
