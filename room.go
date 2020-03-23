package main

type Room struct {
	clients map[*Client]bool
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

// maybe a hub maker?
func newRoom() *Room {
	return &Room{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (room *Room) run() {
	for {
		select {

		// 注册客户端
		case client := <-room.register:
			room.clients[client] = true

		// 注销客户端
		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send)
			}

		// 接受同一房间内客户端发来的广播消息并散布至房间内各个客户端
		case message := <-room.broadcast:
			for client := range room.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(room.clients, client)
				}
			}

		}
	}

}
