package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/handler"
	"github.com/MattBrs/OcelotMDM/internal/api/interceptor"
	"github.com/MattBrs/OcelotMDM/internal/device"
	"github.com/MattBrs/OcelotMDM/internal/storage"
	"github.com/MattBrs/OcelotMDM/internal/token"
	"github.com/MattBrs/OcelotMDM/internal/user"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	userHandler   *api.UserHandler
	tokenHandler  *api.TokenHandler
	deviceHandler *api.DeviceHandler
}

func initAdminUser(userService *user.Service) {
	username := "admin"
	pwd := os.Getenv("ADMIN_PASSWORD")

	if len(pwd) == 0 {
		return
	}

	admin := true
	users, err := userService.QueryUsers(context.TODO(), user.UserFilter{
		Admin: &admin,
	})
	if err != nil {
		fmt.Println("could not create admin user because:", err.Error())
		return
	}

	if len(users) > 0 {
		fmt.Println("there are other admin users in the system")
		return
	}

	adminUser := user.User{
		Username:  username,
		Password:  pwd,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Enabled:   true,
		UpdatedBy: primitive.NewObjectID(),
		Admin:     true,
	}

	err = userService.CreateNewUser(context.TODO(), &adminUser)
	if err != nil {
		fmt.Println("An error occurred while creating admin user")
	}
}

func initMongoConn() *storage.MongoConnection {
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

	if err = mongoConn.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Pinged mongoDb instance successfully")

	return &mongoConn
}

func setGinRoutes(router *gin.Engine, handlers Handlers, authIntr *interceptor.Interceptor) {
	router.POST("/devices", handlers.deviceHandler.AddNewDevice)
	router.GET("/devices", handlers.deviceHandler.ListDevices)
	router.POST("/devices/updateAddress", authIntr.CheckAuth, handlers.deviceHandler.UpdateDeviceAddress)

	router.POST("/token/generate", authIntr.CheckAuth, handlers.tokenHandler.RequestToken)

	router.POST("/user/create", handlers.userHandler.CreateUser)
	router.POST("/user/login", handlers.userHandler.Login)
	router.POST(
		"/user/update/enabled",
		authIntr.CheckAuth,
		handlers.userHandler.EnableUser,
	)

}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mongoConn := initMongoConn()
	defer func() {
		if err = mongoConn.CloseMongoConnection(); err != nil {
			panic(err)
		}
	}()

	userCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "users")
	tokenCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "tokens")
	deviceCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "devices")

	userRepo := user.NewMongoRepository(userCol)
	tokenRepo := token.NewMongoRepository(tokenCol)
	deviceRepo := device.NewMongoRepository(deviceCol)

	userService := user.NewService(userRepo)
	tokenService := token.NewService(tokenRepo)
	deviceService := device.NewService(deviceRepo, tokenService)

	handlers := Handlers{
		userHandler:   api.NewUserHandler(userService),
		tokenHandler:  api.NewTokenHandler(tokenService),
		deviceHandler: api.NewDeviceHandler(deviceService),
	}

	authInterceptor := interceptor.NewAuthInterceptor(userService)

	initAdminUser(userService)

	router := gin.Default()
	setGinRoutes(router, handlers, authInterceptor)

	err = router.Run(":8080") // will expose this later with nginx
	if err != nil {
		panic(err)
	}
}
