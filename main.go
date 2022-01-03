package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	go h.run()

	router := gin.New()
	router.LoadHTMLFiles("index.html")

	router.GET("/list/:listId", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	router.GET("/ws/:listId", func(c *gin.Context) {
		listId := c.Param("listId")
		serveWs(c.Writer, c.Request, listId)
	})

	router.Run("0.0.0.0:8080")
}
