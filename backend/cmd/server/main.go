package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
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
	mongoUser := os.Getenv("MONGO_USERNAME")
	mongoPassword := os.Getenv("MONGO_PASSWORD")
	if mongoUser == "" || mongoPassword == "" {
		panic(errors.New("mongo credentials are not set"))
	}

	mongoConnectionStr := fmt.Sprintf("mongodb+srv://%s:%s@ocelotmdm.oy5pj9q.mongodb.net/?retryWrites=true&w=majority&appName=OcelotMDM", mongoUser, mongoPassword)

	serverApi := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoConnectionStr).SetServerAPIOptions(serverApi)

	mongoClient, err := mongo.Connect(opts)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = mongoClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err = mongoClient.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("Pinged mongoDb instance successfully")

	router := gin.Default()
	router.GET("/ping", ping) // just for testing if everything works, for now :)
	router.Run("localhost:8080")
}
