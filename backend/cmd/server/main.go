package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/MattBrs/OcelotMDM/internal/storage"
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
	mongoConf := storage.DbConfig{
		Username:   os.Getenv("MONGO_USERNAME"),
		Password:   os.Getenv("MONGO_PASSWORD"),
		AppName:    "OcelotMDM",
		ClusterURL: os.Getenv("MONGO_URL"),
	}

	mongoConn, err := storage.NewMongoConnection(mongoConf)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoConn.CloseMongoConnection(); err != nil {
			panic(err)
		}
	}()

	if err = mongoConn.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Pinged mongoDb instance successfully")

	router := gin.Default()
	router.GET("/ping", ping) // just for testing if everything works, for now :)
	router.Run("localhost:8080") // will expose this later with nginx
}
