package ws

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type Message struct {
	Event string `json:"event"`
	Data  string `json:"data"`
}

type Websocket struct {
	// per client represent map[string]*Client type, each client is provided with key
	Clients map[string]*Client

	Broadcast chan []byte

	Register chan *Client

	UnRegister chan *Client
}

// NewWebsocket returns new Websocket
func NewWebsocket() *Websocket {
	return &Websocket{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan []byte),
		Register:   make(chan *Client),
		UnRegister: make(chan *Client),
	}
}

func (ws *Websocket) Run() {
	for {
		select {
		// handle register client case
		case client := <-ws.Register:
			ws.Clients[client.ClientId] = client

			p := Message{
				Event: "rand",
				Data:  "PONG",
			}

			// marshal packet and send in to the channel
			symbolByte, err := json.Marshal(p)
			if err != nil {
				log.Println(err)
				return
			}

			client.Send <- symbolByte

		// unregister client case
		case client := <-ws.UnRegister:
			if _, ok := ws.Clients[client.ClientId]; ok {
				// delete client
				close(client.Send)
				delete(ws.Clients, client.ClientId)
			}

		// handle case to receiving broadcast
		case message := <-ws.Broadcast:
			for _, client := range ws.Clients {
				select {
				case client.Send <- message:
					fmt.Println("Broadcasting client.send")
				default:
					close(client.Send)
					delete(ws.Clients, client.ClientId)
				}
			}
		}
	}
}

func (ws *Websocket) ServeWs(c *gin.Context) {
	query := c.Query("key")
	if query == "" {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "key is required"})
		return
	}

	py, err := ParseJWTClaims(query)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Println(err)
		return
	}

	conn, err := ConnectionUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		fmt.Println(err)

		return
	}

	// initialize websocket client
	client := &Client{
		Ws:       ws,
		Conn:     conn,
		ClientId: strconv.Itoa(int(py.UserId)),
		Send:     make(chan []byte, 256),
	}

	// register initialized client
	client.Ws.Register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// another goroutines.
	go client.WritePump()
	go client.ReadPump()
}
