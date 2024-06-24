package ws

import (
	"encoding/json"
	"fmt"
	"github.com/Nicolas-ggd/ch-mod/internal/db/models/request"
	"github.com/Nicolas-ggd/ch-mod/pkg/services"
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

	ChatHandler services.IChatService
}

// NewWebsocket returns new Websocket
func NewWebsocket(authService services.IChatService) *Websocket {
	return &Websocket{
		Clients:     make(map[string]*Client),
		Broadcast:   make(chan []byte),
		Register:    make(chan *Client),
		UnRegister:  make(chan *Client),
		ChatHandler: authService,
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

// SendEvent function send events to the client
func (ws *Websocket) SendEvent(clients []string, data []byte) {
	var cl *Client
	var m request.WsChatRequest

	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Println(err)
	}

	value, err := json.Marshal(&m)
	if err != nil {
		log.Printf("Can't marshal action data")
		return
	}

	for _, client := range clients {
		c, ok := ws.Clients[client]
		if !ok {
			log.Printf("Client with ID %s not found", client)
			continue // Continue process
		}

		cl = c
	}

	_, err = ws.ChatHandler.Create(&m)
	if err != nil {
		log.Println(err)
	}

	// Check if the Send channel is initialized
	if cl.Send == nil {
		log.Printf("Send channel not initialized for client with ID %v", cl)
		return // Exit the function without sending data
	}

	// Send data to the client
	cl.Send <- value
}

// BroadcastEvent function send events in broadcast
func (ws *Websocket) BroadcastEvent(data []byte) {
	var m request.WsChatRequest

	err := json.Unmarshal(data, &m)
	if err != nil {
		log.Println(err)
	}

	value, err := json.Marshal(&m)
	if err != nil {
		log.Printf("Can't marshal broadcast data")
		return
	}

	// Send data to the Broadcast channel
	ws.Broadcast <- value
}
