package command_queue

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/domain/command"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandQueueService struct {
	ticker         *time.Ticker
	doneChannel    chan bool
	mqttClient     *ocelot_mqtt.MqttClient
	commandService *command.Service
	ctx            context.Context
}

func NewService(
	context context.Context,
	messageHandler *ocelot_mqtt.MqttClient,
	cmdService *command.Service,
	tickerInterval time.Duration,
) *CommandQueueService {
	service := CommandQueueService{
		ticker:         time.NewTicker(tickerInterval),
		doneChannel:    make(chan bool),
		mqttClient:     messageHandler,
		commandService: cmdService,
		ctx:            context,
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
			case <-service.ticker.C:
				onFetch(service)
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
				onAckResponse(service, &msg)
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

func fetchWaitingCmds(s *CommandQueueService) ([]*command.Command, error) {
	commands, err := s.commandService.ListCommands(
		s.ctx, command.CommandFilter{
			Status: &command.WAITING,
		},
	)

	if err != nil {
		return nil, err
	}

	return commands, nil
}

func enqueueWaitingCmds(
	s *CommandQueueService,
	cmds []*command.Command,
) (*primitive.ObjectID, error) {
	queueID := primitive.NewObjectID()
	err := s.commandService.EnqueueMany(s.ctx, cmds, queueID)
	if err != nil {
		return nil, err
	}

	return &queueID, nil
}

func onFetch(s *CommandQueueService) {
	commands, err := fetchWaitingCmds(s)
	if err != nil {
		fmt.Println("error while fetching the commands: ", err.Error())
		return
	}

	if len(commands) == 0 {
		fmt.Println("no commands to enqueue")
		return
	}

	queueID, err := enqueueWaitingCmds(s, commands)
	if err != nil {
		fmt.Println("commands were not enqueued because: ", err.Error())
		return
	}

	fmt.Println("commands have beed enqueued with ID: ", queueID.Hex())

	for i := range commands {
		// TODO(matteobrusarosco):
		// sent commands should contain ID, cmd_type and payload
		err = s.mqttClient.Publish(
			commands[i].Id.Hex()+commands[i].CommandActionName,
			commands[i].DeviceName+"/cmd",
			1,
		)

		if err != nil {
			_ = s.commandService.UpdateStatus(
				s.ctx,
				commands[i].Id.Hex(),
				command.ERRORED,
				fmt.Sprintf(
					"could not send to device because: %s",
					err.Error(),
				),
			)
		}
	}
}

func onAckResponse(s *CommandQueueService, msg *ocelot_mqtt.ChanMessage) {
	splittedTopic := strings.Split(msg.Topic, "/")
	if len(splittedTopic) != 2 {
		fmt.Println("is not ack from device")
		return
	}

	// on ack topic there are only IDs of acked commands
	err := s.commandService.UpdateStatus(
		s.ctx,
		string(msg.Payload),
		command.ACKED,
		"",
	)

	if err != nil {
		fmt.Println("could not update command status: ", err.Error())
		return
	}

	fmt.Println(" on reponse, received msg from topic: ", msg.Topic)
}
