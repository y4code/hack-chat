package main

// Hub 维持了一些 激活的客户端 和 发往客户端的信息
type Hub struct {
	// 注册的客户端
	clients map[*Client]bool

	// 从客户端发来的信息
	broadcast chan []byte

	// 从客户端注册的请求
	register chan *Client

	// 从客户端注销的请求
	unregister chan *Client
}

// maybe a hub maker?
func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {

		// 注册客户端
		case client := <-h.register:
			h.clients[client] = true

		// 注销客户端
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		// 接受客户端发来的广播消息并散布至各个客户端
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}

		}
	}

}
