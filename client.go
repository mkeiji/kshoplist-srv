package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

type KmsgType string

const (
	REQ  KmsgType = "Request"
	RESP KmsgType = "Response"
)

type Actions string

const (
	PUT    Actions = "PUT"
	POST   Actions = "POST"
	DELETE Actions = "DELETE"
)

type Kmsg struct {
	Type   KmsgType `json:"type"`
	Action Actions  `json:"action"`
	Items  []Item   `json:"items"`
}

type Item struct {
	Id      int    `json:"id"`
	StoreId int    `json:"storeId"`
	Name    string `json:"name"`
}

var mockList []Item = []Item{
	{Id: 1, StoreId: 1, Name: "Hello"},
	{Id: 2, StoreId: 2, Name: "World"},
}

var mockId int = 3

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type connection struct {
	// The websocket connection.
	ws *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
func (s subscription) readPump() {
	c := s.conn
	defer func() {
		h.unregister <- s
		c.ws.Close()
	}()
	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		// Todo: save to DB here
		var newMsg []byte
		var kmsg Kmsg
		kmsgStr := string(msg)
		kerr := json.Unmarshal([]byte(kmsgStr), &kmsg)
		if kerr != nil {
			fmt.Println(kerr)
		}
		fmt.Printf("\nreading: %+v\n", kmsg)
		if kmsg.Type == REQ && kmsg.Action == POST {
			item := kmsg.Items[0]
			item.Id = mockId
			mockId++

			mockList = append(mockList, item)
			fmt.Printf("\npostOnDb: %+v\n", mockList)
			response := Kmsg{
				Type:  RESP,
				Items: mockList,
			}
			newMsg, _ = json.Marshal(response)
		}
		if kmsg.Type == REQ && kmsg.Action == PUT {
			var updatedMockList []Item
			newItem := kmsg.Items[0]
			for _, i := range mockList {
				if i.Id == newItem.Id {
					updatedMockList = append(updatedMockList, Item{
						Id:      i.Id,
						StoreId: i.StoreId,
						Name:    newItem.Name,
					})
				} else {
					updatedMockList = append(updatedMockList, i)
				}
			}
			mockList = updatedMockList

			fmt.Printf("\nupdateDb: %+v\n", updatedMockList)
			response := Kmsg{
				Type:  RESP,
				Items: updatedMockList,
			}
			newMsg, _ = json.Marshal(response)
		}
		if kmsg.Type == REQ && kmsg.Action == DELETE {
			var updatedMockList []Item
			newItem := kmsg.Items[0]
			for index, i := range mockList {
				if i.Id == newItem.Id {
					updatedMockList = append(
						updatedMockList[:index],
						updatedMockList[index:]...,
					)
					fmt.Printf("\nHERE: %v\n", updatedMockList)
				} else {
					updatedMockList = append(updatedMockList, i)
				}
			}
			mockList = updatedMockList

			fmt.Printf("\nupdateDb: %+v\n", updatedMockList)
			response := Kmsg{
				Type:  RESP,
				Items: updatedMockList,
			}
			newMsg, _ = json.Marshal(response)
		}

		m := message{newMsg, s.list}
		h.broadcast <- m
	}
}

// write writes a message with the given message type and payload.
func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, payload)
}

// writePump pumps messages from the hub to the websocket connection.
func (s *subscription) writePump() {
	c := s.conn
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func serveWs(w http.ResponseWriter, r *http.Request, listId string) {
	fmt.Printf("\nnew client connected on list #%v\n", listId)
	upgrader.CheckOrigin = func(r *http.Request) bool {
		fmt.Printf("request origin: %v", r.Header["Origin"])
		return true
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err.Error())
		return
	}
	c := &connection{send: make(chan []byte, 256), ws: ws}
	s := subscription{c, listId}
	h.register <- s

	// Read from Db and send current list state
	welcomeMsg := Kmsg{
		Type:  RESP,
		Items: mockList,
	}
	msg, _ := json.Marshal(welcomeMsg)
	c.write(websocket.TextMessage, msg)

	go s.writePump()
	go s.readPump()
}
