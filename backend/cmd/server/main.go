package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/api/handler"
	"github.com/MattBrs/OcelotMDM/internal/api/interceptor"
	"github.com/MattBrs/OcelotMDM/internal/domain/binary"
	"github.com/MattBrs/OcelotMDM/internal/domain/command"
	"github.com/MattBrs/OcelotMDM/internal/domain/command_action"
	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"github.com/MattBrs/OcelotMDM/internal/domain/file_repository"
	"github.com/MattBrs/OcelotMDM/internal/domain/logs"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
	"github.com/MattBrs/OcelotMDM/internal/domain/service/command_queue"
	"github.com/MattBrs/OcelotMDM/internal/domain/service/logs_handler"
	"github.com/MattBrs/OcelotMDM/internal/domain/service/uptime"
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
	binaryHandler        *api.BinaryHandler
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

	router.POST(
		"/binary/add",
		authIntr.CheckAuth,
		handlers.binaryHandler.AddNewBinary,
	)
	router.GET(
		"/binary/get",
		handlers.binaryHandler.GetBinary,
	)
}

func newMqttClient() *ocelot_mqtt.MqttClient {
	mqttHost := os.Getenv("MQTT_HOST")
	mqttPort, err := strconv.Atoi(os.Getenv("MQTT_PORT"))
	if err != nil {
		fmt.Println("mqtt port is not a number")
		panic(1)
	}

	mqttClient := ocelot_mqtt.NewMqttClient(mqttHost, uint(mqttPort))

	err = mqttClient.Connect()
	if err != nil {
		fmt.Println("unable to establish mqtt connection")
		panic(1)
	}

	return mqttClient
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

	mqttClient := newMqttClient()
	defer mqttClient.Close()

	logCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "logs")
	userCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "users")
	tokenCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "tokens")
	binaryCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "binaries")
	deviceCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "devices")
	commandCol := mongoConn.GetCollection(os.Getenv("DBNAME"), "commands")
	commandActionCol := mongoConn.GetCollection(
		os.Getenv("DBNAME"), "command_actions",
	)

	s3Repo := file_repository.NewS3Repository(
		context.Background(),
		os.Getenv("SPACES_KEY"),
		os.Getenv("SPACES_SECRET"),
		os.Getenv("SPACES_ENDPOINT"),
		os.Getenv("SPACES_BUCKET"),
		os.Getenv("SPACES_REGION"),
	)

	if s3Repo == nil {
		fmt.Println("s3 client is nil")
		os.Exit(1)
	}

	fmt.Println("s3 client created")

	logRepo := logs.NewMongoRepository(logCol)
	userRepo := user.NewMongoRepository(userCol)
	tokenRepo := token.NewMongoRepository(tokenCol)
	binaryRepo := binary.NewMongoRepository(binaryCol)
	deviceRepo := device.NewMongoRepository(deviceCol)
	commandRepo := command.NewMongoRepository(commandCol)
	commandActionRepo := command_action.NewMongoCommandActionRepository(
		commandActionCol,
	)

	logService := logs.NewService(logRepo, s3Repo)
	userService := user.NewService(userRepo)
	tokenService := token.NewService(tokenRepo)
	vpnService := vpn.NewService("http://vpn_api:8080")
	commandActionService := command_action.NewService(commandActionRepo)
	binaryService := binary.NewService(s3Repo, binaryRepo, tokenService)
	deviceService := device.NewService(
		deviceRepo,
		tokenService,
		vpnService,
		mqttClient,
	)
	commandService := command.NewService(
		commandRepo,
		deviceService,
		commandActionService,
	)

	handlers := Handlers{
		userHandler:    api.NewUserHandler(userService),
		tokenHandler:   api.NewTokenHandler(tokenService),
		binaryHandler:  api.NewBinaryHandler(binaryService),
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

	commandQueueService := command_queue.NewService(
		context.Background(),
		mqttClient,
		commandService,
		tokenService,
		time.Second*10,
	)
	commandQueueService.Start()

	uptimeService := uptime.NewService(
		context.Background(),
		mqttClient,
		deviceService,
	)
	uptimeService.Start()

	logHandlerService := logs_handler.NewService(
		context.Background(),
		mqttClient,
		logService,
	)
	logHandlerService.Start()

	defer logHandlerService.Stop()
	defer uptimeService.Stop()
	defer commandQueueService.Stop()

	err = router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
