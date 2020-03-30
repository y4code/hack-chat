package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 允许写入一条信息到peer的时间
	writeWait = 10 * time.Second

	// 允许从下一个从peer传来的pong 信息的读取时间
	pongWait = 60 * time.Second

	// 发送到 peer 的 ping在这个时间段之内，必须小于pongWait
	pingPeriod = (pongWait * 9) / 10

	// 从peer传来最大消息大小
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 是在 websocket 连接和 Hub 之间的中间人
type Client struct {
	hub *Hub

	nick string

	//use room in case of confuse with channel which as go type
	room *Room

	// websocket 连接
	conn *websocket.Conn

	// 即将发送出去消息的缓冲 channel
	send chan []byte
}

func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	// 设置读取限制大小
	c.conn.SetReadLimit(maxMessageSize)

	// 设置读取DeadLine
	c.conn.SetReadDeadline(time.Now().Add(pongWait))

	// 设置 心跳 pong 回调处理，在回调中，再以现在时间点 重新设置DeadLine，
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// 读取ws消息
		_, message, err := c.conn.ReadMessage()
		fmt.Println(string(message))

		// save nick name to return a json containing property "nick"
		var msg Message
		json.Unmarshal([]byte(string(message)), &msg)
		if msg.Cmd == "join" && len(msg.Nick) > 0 {
			c.nick = msg.Nick
		}

		if err != nil {
			// ws 1001 1006 是否为预期中的错误，不是的话，print
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.room.broadcast <- message

	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			fmt.Println(message)

			addPropertyToMessage(c.nick, &message)

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}

		}

	}
}

func addPropertyToMessage(property string, message *[]byte) {

	var joinMsg JoinMessage
	var chatMsg ChatMessage
	json.Unmarshal([]byte(string(*message)), &joinMsg)
	json.Unmarshal([]byte(string(*message)), &chatMsg)
	joinMsg.Nick = property
	chatMsg.Nick = property
	// TODO 需要优化
	if joinMsg.Cmd == "chat" {
		*message, _ = json.Marshal(chatMsg)
	} else {
		*message, _ = json.Marshal(joinMsg)
	}
	fmt.Println("msg 是", string(*message))
}

type CanNotGetRoomNumber string

func (e CanNotGetRoomNumber) Error() string {
	return "Can not get room number"
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	//升级连接
	conn, _ := upgrader.Upgrade(w, r, nil)
	//if err != nil {
	//	log.Println(err)
	//	return
	//}

	//新建一个 client
	client := &Client{
		hub:  hub,
		nick: "",
		room: getRoom(hub, r),
		conn: conn,
		send: make(chan []byte, 256),
	}

	//client.hub.rooms[r.URL.Path].register <- client
	//房间入住客人，这里上下两行皆可
	client.room.register <- client

	// 在新的协程中完成所有工作 以允许(调用者)引用一些内存
	go client.writePump()
	go client.readPump()
}

func getRoom(hub *Hub, r *http.Request) *Room {
	roomNumber, e := getRoomNumber(hub, r)
	if e != nil {
		fmt.Println(e)
	}

	// 有房的话，返回房间；没房的话，新建房间
	if _, ok := hub.rooms[roomNumber]; !ok {
		nRoom := newRoom()
		//第一次开房先要初始化房间
		go nRoom.run()
		// TODO set room in gouroutine maybe a better perfomance
		hub.rooms[roomNumber] = nRoom
		return nRoom
	} else {
		//有房的话，给房卡
		return hub.rooms[roomNumber]
	}
}

func getRoomNumber(hub *Hub, r *http.Request) (string, error) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return "", CanNotGetRoomNumber("")
	}
	roomNumber := r.Form.Get("channel")
	fmt.Println(roomNumber)
	return roomNumber, nil
}
