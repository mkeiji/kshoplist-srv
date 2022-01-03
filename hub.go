package main

type message struct {
	data []byte
	list string
}

type subscription struct {
	conn *connection
	list string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	lists map[string]map[*connection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	lists:      make(map[string]map[*connection]bool),
}

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.lists[s.list]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.lists[s.list] = connections
			}
			h.lists[s.list][s.conn] = true
		case s := <-h.unregister:
			connections := h.lists[s.list]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.lists, s.list)
					}
				}
			}
		case m := <-h.broadcast:
			connections := h.lists[m.list]
			for c := range connections {
				select {
				case c.send <- m.data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.lists, m.list)
					}
				}
			}
		}
	}
}
