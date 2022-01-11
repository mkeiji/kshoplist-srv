package services

import (
	"encoding/json"
	"fmt"
	"kshoplistSrv/constants"
	e "kshoplistSrv/enums"
	m "kshoplistSrv/models"
	r "kshoplistSrv/repository"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// consider using a DB
var cacheList []m.Item = []m.Item{}
var nextId int = 1

type Subscription struct {
	Conn           *m.Connection
	ListId         string
	ItemRepository r.ItemRepository
}

func (this Subscription) OnConnInit(c *m.Connection) {
	// get all items on first connection
	welcomeMsg := m.Kdata{
		Type:  e.RESP,
		Items: this.ItemRepository.GetAll(),
	}
	msg, _ := json.Marshal(welcomeMsg)
	c.Write(websocket.TextMessage, msg)
}

// readPump pumps messages from the websocket connection to the hub.
func (this Subscription) ReadPump(h Hub) {
	c := this.Conn
	defer func() {
		h.Unregister <- this
		c.Ws.Close()
	}()
	c.Ws.SetReadLimit(constants.MaxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(constants.PongWait))
	c.Ws.SetPongHandler(func(string) error {
		c.Ws.SetReadDeadline(time.Now().Add(constants.PongWait))
		return nil
	})
	for {
		_, rawMsg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("INFO: %v", err)
			}
			break
		}
		newMsg := this.handleRawMsg(rawMsg)
		m := m.Message{Data: newMsg, List: this.ListId}
		h.Broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (this *Subscription) WritePump() {
	c := this.Conn
	ticker := time.NewTicker(constants.PingPeriod)
	defer func() {
		ticker.Stop()
		c.Ws.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Write(websocket.TextMessage, message); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.Write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

/* PRIVATE
--------------------------------------------------------------------------------*/
// request handler
func (this Subscription) handleRawMsg(msg []byte) []byte {
	var newMsg []byte
	var kreq m.Kdata

	kmsgStr := string(msg)
	kerr := json.Unmarshal([]byte(kmsgStr), &kreq)
	if kerr != nil {
		fmt.Println(kerr)
	}
	if kreq.Type == e.REQ {
		switch kreq.Action {
		case e.POST:
			fmt.Printf("\n[POST]: %+v\n", kreq)
			newMsg = this.handlePost(kreq)
			break
		case e.PUT:
			fmt.Printf("\n[PUT]: %+v\n", kreq)
			newMsg = this.handlePut(kreq)
			break
		case e.DELETE:
			fmt.Printf("\n[DELETE]: %+v\n", kreq)
			newMsg = this.handleDelete(kreq)
			break
		}
	}
	return newMsg
}

func (this Subscription) handlePost(kreq m.Kdata) []byte {
	item := kreq.Items[0]
	this.ItemRepository.Post(item)
	return this.getDefaultResponse()
}

func (this Subscription) handlePut(kreq m.Kdata) []byte {
	newItem := kreq.Items[0]
	this.ItemRepository.Put(newItem)
	return this.getDefaultResponse()
}

func (this Subscription) handleDelete(kreq m.Kdata) []byte {
	itemToDelete := kreq.Items[0]
	this.ItemRepository.Delete(itemToDelete)
	return this.getDefaultResponse()
}

func (this Subscription) getDefaultResponse() []byte {
	response := m.Kdata{
		Type:  e.RESP,
		Items: this.ItemRepository.GetAll(),
	}
	newMsg, _ := json.Marshal(response)
	return newMsg
}
