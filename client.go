package room

import "github.com/gorilla/websocket"

type Client struct {
	socket  *websocket.Conn
	recieve chan []byte
	room    *Room
}

func (c *Client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		c.room.Forward <- msg
	}
}
func (c *Client) write() {
	defer c.socket.Close()
	for msg := range c.recieve {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
