package services

import (
	m "kshoplistSrv/models"
)

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	Lists map[string]map[*m.Connection]bool

	// Inbound messages from the connections.
	Broadcast chan m.Message

	// Register requests from the connections.
	Register chan Subscription

	// Unregister requests from connections.
	Unregister chan Subscription
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Lists[s.ListId]
			if connections == nil {
				connections = make(map[*m.Connection]bool)
				h.Lists[s.ListId] = connections
			}
			h.Lists[s.ListId][s.Conn] = true
		case s := <-h.Unregister:
			connections := h.Lists[s.ListId]
			if connections != nil {
				if _, ok := connections[s.Conn]; ok {
					delete(connections, s.Conn)
					close(s.Conn.Send)
					if len(connections) == 0 {
						delete(h.Lists, s.ListId)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Lists[m.List]
			for c := range connections {
				select {
				case c.Send <- m.Data:
				default:
					close(c.Send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Lists, m.List)
					}
				}
			}
		}
	}
}
