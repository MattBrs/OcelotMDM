package logs

import (
	"context"
	"fmt"

	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
)

type LogService struct {
	mqttClient  *ocelot_mqtt.MqttClient
	doneChannel chan bool
	ctx         context.Context
}

func NewService(ctx context.Context, mqttClient *ocelot_mqtt.MqttClient) *LogService {
	return &LogService{
		mqttClient:  mqttClient,
		doneChannel: make(chan bool),
		ctx:         ctx,
	}
}

func (s *LogService) Start() {
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

func (s *LogService) Stop() {
	s.doneChannel <- true
	close(s.doneChannel)
}

func onMsgReceived(service *LogService, msg *ocelot_mqtt.ChanMessage) {
	fmt.Printf("arrived log of size %d from topic %s\n", len(msg.Payload), msg.Topic)
}
