package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/handler"
	"github.com/MattBrs/OcelotMDM/internal/api/interceptor"
	"github.com/MattBrs/OcelotMDM/internal/domain/command"
	"github.com/MattBrs/OcelotMDM/internal/domain/command_action"
	"github.com/MattBrs/OcelotMDM/internal/domain/command_queue"
	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/MattBrs/OcelotMDM/internal/domain/user"
	"github.com/MattBrs/OcelotMDM/internal/domain/vpn"
	"github.com/MattBrs/OcelotMDM/internal/storage"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Handlers struct {
	userHandler          *api.UserHandler
	tokenHandler         *api.TokenHandler
	deviceHandler        *api.DeviceHandler
	commandHandler       *api.CommandHandler
	commandActionHandler *api.CommandActionHandler
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
	router.GET(
		"/devices",
		authIntr.CheckAuth,
		handlers.deviceHandler.ListDevices,
	)
	router.POST(
		"/devices/updateAddress",
		authIntr.CheckAuth,
		handlers.deviceHandler.UpdateDeviceAddress,
	)

	router.POST(
		"/token/generate",
		authIntr.CheckAuth,
		handlers.tokenHandler.RequestToken,
	)

	router.POST("/user/create", handlers.userHandler.CreateUser)
	router.POST("/user/login", handlers.userHandler.Login)
	router.POST(
		"/user/update/enabled",
		authIntr.CheckAuth,
		handlers.userHandler.EnableUser,
	)

	router.POST(
		"/command/new",
		authIntr.CheckAuth,
		handlers.commandHandler.AddNewCommand,
	)
	router.GET(
		"/command/list",
		authIntr.CheckAuth,
		handlers.commandHandler.ListCommands,
	)
	router.POST(
		"/command/update/status",
		authIntr.CheckAuth,
		handlers.commandHandler.UpdateCommandStatus,
	)

	router.POST(
		"/command_actions/new",
		authIntr.CheckAuth,
		handlers.commandActionHandler.AddNewCommandAction,
	)
	router.GET(
		"/command_actions/list",
		authIntr.CheckAuth,
		handlers.commandActionHandler.ListCommandActions,
	)
	router.POST(
		"/command_actions/delete",
		authIntr.CheckAuth,
		handlers.commandActionHandler.DeleteCommandAction,
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
	commandCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "commands")
	commandActionCol := mongoConn.GetCollection(
		os.Getenv("DBNAME"), "command_actions",
	)

	userRepo := user.NewMongoRepository(userCol)
	tokenRepo := token.NewMongoRepository(tokenCol)
	deviceRepo := device.NewMongoRepository(deviceCol)
	commandRepo := command.NewMongoRepository(commandCol)
	commandActionRepo := command_action.NewMongoCommandActionRepository(
		commandActionCol,
	)

	vpnService := vpn.NewService("http://vpn_api:8080")
	userService := user.NewService(userRepo)
	tokenService := token.NewService(tokenRepo)
	deviceService := device.NewService(deviceRepo, tokenService, vpnService)
	commandActionService := command_action.NewService(commandActionRepo)
	commandService := command.NewService(
		commandRepo,
		deviceService,
		commandActionService,
	)

	handlers := Handlers{
		userHandler:    api.NewUserHandler(userService),
		tokenHandler:   api.NewTokenHandler(tokenService),
		deviceHandler:  api.NewDeviceHandler(deviceService),
		commandHandler: api.NewCommandHandler(commandService),
		commandActionHandler: api.NewCommandActionHandler(
			commandActionService,
		),
	}

	authInterceptor := interceptor.NewAuthInterceptor(userService)

	initAdminUser(userService)

	router := gin.Default()
	setGinRoutes(router, handlers, authInterceptor)

	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort, err := strconv.Atoi(os.Getenv("MQTT_PORT"))
	if err != nil {
		fmt.Println("mqtt port is not a number")
		panic(1)
	}

	pahoClient := ocelot_mqtt.NewMqttClient(
		mqttHost,
		uint(mqttPort),
	)

	err = pahoClient.Connect()
	if err != nil {
		fmt.Println("unable to establish mqtt connection")
		panic(1)
	}

	defer pahoClient.Close()

	err = pahoClient.Subscribe("misty-dew/ack", 0)
	if err != nil {
		fmt.Println("error un subscription for topic test")
	}

	commandQueueService := command_queue.NewService(
		context.Background(),
		pahoClient,
		commandService,
		time.Second*10,
	)
	commandQueueService.Start()

	defer commandQueueService.Stop()

	err = router.Run(":8080") // will expose this later with nginx
	if err != nil {
		panic(err)
	}
}
