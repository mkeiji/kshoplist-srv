package models

import (
	"encoding/json"
	"fmt"
	e "kshoplistSrv/enums"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// consider using a DB
var cacheList []Item = []Item{}
var nextId int = 1

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

type Subscription struct {
	Conn   *Connection
	ListId string
}

func (this Subscription) OnConnInit(c *Connection) {
	// Read from cache and send current list state
	welcomeMsg := Kdata{
		Type:  e.RESP,
		Items: cacheList,
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
	c.Ws.SetReadLimit(maxMessageSize)
	c.Ws.SetReadDeadline(time.Now().Add(pongWait))
	c.Ws.SetPongHandler(func(string) error { c.Ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, rawMsg, err := c.Ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("INFO: %v", err)
			}
			break
		}
		newMsg := this.handleRawMsg(rawMsg)
		m := Message{newMsg, this.ListId}
		h.Broadcast <- m
	}
}

// writePump pumps messages from the hub to the websocket connection.
func (this *Subscription) WritePump() {
	c := this.Conn
	ticker := time.NewTicker(pingPeriod)
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
	// Todo: save to DB here
	var newMsg []byte
	var kreq Kdata

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

func (this Subscription) handlePost(kreq Kdata) []byte {
	item := kreq.Items[0]
	item.Id = nextId
	nextId++

	cacheList = append(cacheList, item)
	response := Kdata{
		Type:  e.RESP,
		Items: cacheList,
	}
	newMsg, _ := json.Marshal(response)
	return newMsg
}

func (this Subscription) handlePut(kreq Kdata) []byte {
	var updatedMockList []Item
	newItem := kreq.Items[0]
	for _, i := range cacheList {
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
	cacheList = updatedMockList

	response := Kdata{
		Type:  e.RESP,
		Items: updatedMockList,
	}
	newMsg, _ := json.Marshal(response)
	return newMsg
}

func (this Subscription) handleDelete(kreq Kdata) []byte {
	var updatedMockList []Item
	newItem := kreq.Items[0]
	for index, i := range cacheList {
		if i.Id == newItem.Id {
			updatedMockList = append(
				updatedMockList[:index],
				updatedMockList[index:]...,
			)
		} else {
			updatedMockList = append(updatedMockList, i)
		}
	}
	cacheList = updatedMockList

	response := Kdata{
		Type:  e.RESP,
		Items: updatedMockList,
	}
	newMsg, _ := json.Marshal(response)
	return newMsg
}
