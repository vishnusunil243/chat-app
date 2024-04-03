package room

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	Clients map[*Client]bool
	Join    chan *Client
	Leave   chan *Client
	Forward chan []byte
}

func NewRoom() *Room {
	return &Room{
		Clients: make(map[*Client]bool),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Forward: make(chan []byte),
	}
}
func (r *Room) Run() {
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			delete(r.Clients, client)
			close(client.recieve)
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.recieve <- msg
			}
		}
	}
}

const (
	socketBufferSize  = 1824
	messageBufferSize = 236
)

var upgrader = websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *Room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("read error", err)
		return
	}
	client := &Client{
		socket:  socket,
		recieve: make(chan []byte, messageBufferSize),
		room:    r,
	}
	r.Join <- client
	defer func() { r.Leave <- client }()
	go client.write()
	client.read()
}
