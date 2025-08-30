package command_queue

import (
	"fmt"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
)

type CommandQueueService struct {
	ticker      *time.Ticker
	doneChannel chan bool
	mqttClient  *ocelot_mqtt.MqttClient
}

func NewService(
	messageHandler *ocelot_mqtt.MqttClient,
	tickerInterval time.Duration,
) *CommandQueueService {
	service := CommandQueueService{
		ticker:      time.NewTicker(tickerInterval),
		doneChannel: make(chan bool),
		mqttClient:  messageHandler,
	}

	return &service
}

func (s *CommandQueueService) Start() {
	startSender(s)
	startReceiver(s)
}

func startSender(service *CommandQueueService) {
	go func() {
		for {
			select {
			case <-service.doneChannel:
				return
			case t := <-service.ticker.C:
				// read waiting commands from DB and enqueue them on mqtt
				fmt.Println("ticker ticked at ", t)
			}
		}
	}()
}

func startReceiver(service *CommandQueueService) {
	go func() {
		for {
			select {
			case <-service.doneChannel:
				return
			case msg := <-service.mqttClient.AckMessages:
				// handle msg
				fmt.Println("received msg from topic: ", msg.Topic)
			}
		}
	}()

}

func (service *CommandQueueService) Stop() {
	service.ticker.Stop()
	service.doneChannel <- true

	close(service.doneChannel)

	fmt.Println("queue service stopped successfully")
}
