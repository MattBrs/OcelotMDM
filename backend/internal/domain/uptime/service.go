package uptime

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/MattBrs/OcelotMDM/internal/domain/device"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
)

type UptimeService struct {
	mqttClient    *ocelot_mqtt.MqttClient
	deviceService *device.Service
	doneChannel   chan bool
	ctx           context.Context
}

func NewService(
	ctx context.Context,
	mqttClient *ocelot_mqtt.MqttClient,
	deviceService *device.Service,
) *UptimeService {
	return &UptimeService{
		mqttClient:    mqttClient,
		deviceService: deviceService,
		doneChannel:   make(chan bool),
		ctx:           ctx,
	}
}

func (s *UptimeService) Start() {
	go func() {
		for {
			select {
			case <-s.doneChannel:
				return
			case msg := <-s.mqttClient.UptimeMessages:
				onMsgReceived(s, &msg)
			}
		}
	}()
}

func (s *UptimeService) Stop() {
	s.doneChannel <- true
	close(s.doneChannel)
}

func onMsgReceived(service *UptimeService, msg *ocelot_mqtt.ChanMessage) {
	topicParts := strings.Split(msg.Topic, "/")
	msgParts := strings.Split(string(msg.Payload), " ")

	if len(topicParts) != 2 || strings.Compare(topicParts[1], "online") != 0 {
		fmt.Println("topic not the uptime one: ", msg.Topic)
		return
	}

	if len(msgParts) != 2 {
		fmt.Println("wrong number of params")
		return
	}

	epoch, err := strconv.ParseInt(msgParts[0], 10, 64)
	ip := msgParts[1]

	if err != nil {
		fmt.Println("second param is not epoch int64: ", msgParts[0])
		return
	}

	fmt.Println("try update uptime status")
	err = service.deviceService.UpdateUpStatus(service.ctx, topicParts[0], ip, epoch)
	if err != nil {
		fmt.Println("error while updating ip and lastSeen")
	}
}
