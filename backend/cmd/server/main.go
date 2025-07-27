package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// gin test structure
type hello struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// gin test function
func ping(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, hello{Title: "pong", Message: "Hi dude!"})
}

func main() {
	// TODO: initialize mongoDB connection, initialize services and add
	// proper api calls

	router := gin.Default()
	router.GET("/ping", ping) // just for testing if everything works, for now :)
	router.Run("localhost:8080")
}
