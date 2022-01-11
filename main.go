package main

import (
	"fmt"
	m "kshoplistSrv/models"
	"kshoplistSrv/repository"
	s "kshoplistSrv/services"
	"log"
	"net/http"

	appdb "kshoplistSrv/database"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var appDb appdb.AppDb

var h = s.Hub{
	Broadcast:  make(chan m.Message),
	Register:   make(chan s.Subscription),
	Unregister: make(chan s.Subscription),
	Lists:      make(map[string]map[*m.Connection]bool),
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	appDb.Init()
	defer appdb.Db.Close()

	go h.Run()

	router := gin.New()
	router.LoadHTMLFiles("templates/index.html")

	router.GET("/list/:listId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:listId", func(c *gin.Context) {
		listId := c.Param("listId")
		serveWs(c.Writer, c.Request, listId)
	})

	router.Run("0.0.0.0:8081")
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
	c := &m.Connection{Send: make(chan []byte, 256), Ws: ws}
	s := s.Subscription{Conn: c, ListId: listId, ItemRepository: repository.NewItemRepository()}
	h.Register <- s

	s.OnConnInit(c)
	go s.WritePump()
	go s.ReadPump(h)
}
