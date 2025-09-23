package logs_handler

import (
	"context"
	"fmt"
	"strings"

	"github.com/MattBrs/OcelotMDM/internal/domain/logs"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
)

type LogHandlerService struct {
	mqttClient  *ocelot_mqtt.MqttClient
	logsService *logs.Service
	doneChannel chan bool
	ctx         context.Context
}

func NewService(ctx context.Context, mqttClient *ocelot_mqtt.MqttClient, logsService *logs.Service) *LogHandlerService {
	return &LogHandlerService{
		mqttClient:  mqttClient,
		doneChannel: make(chan bool),
		ctx:         ctx,
		logsService: logsService,
	}
}

func (s *LogHandlerService) Start() {
	go func() {
		for {
			select {
			case <-s.doneChannel:
				return
			case msg := <-s.mqttClient.LogMessages:
				onMsgReceived(s, &msg)
			}
		}
	}()
}

func (s *LogHandlerService) Stop() {
	s.doneChannel <- true
	close(s.doneChannel)
}

func onMsgReceived(service *LogHandlerService, msg *ocelot_mqtt.ChanMessage) {
	topicParts := strings.Split(msg.Topic, "/")

	if len(topicParts) != 2 || topicParts[1] != "logs" {
		fmt.Println("topic is not logs")
		return
	}

	fmt.Printf("arrived log of size %d from device %s\n", len(msg.Payload), topicParts[0])
	err := service.logsService.AddLog(service.ctx, topicParts[0], msg.Payload)

	if err != nil {
		fmt.Printf("error while inserting new log from device %s\n", topicParts[0])
	}
}
