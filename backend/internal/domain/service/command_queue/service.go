package command_queue

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/MattBrs/OcelotMDM/internal/domain/command"
	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/ocelot_mqtt"
	"github.com/MattBrs/OcelotMDM/internal/domain/token"
	"github.com/vmihailenco/msgpack/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CommandQueueService struct {
	ticker         *time.Ticker
	doneChannel    chan bool
	mqttClient     *ocelot_mqtt.MqttClient
	commandService *command.Service
	tokenService   *token.Service
	ctx            context.Context
}

func NewService(
	context context.Context,
	messageHandler *ocelot_mqtt.MqttClient,
	cmdService *command.Service,
	tokenService *token.Service,
	tickerInterval time.Duration,
) *CommandQueueService {
	service := CommandQueueService{
		ticker:         time.NewTicker(tickerInterval),
		doneChannel:    make(chan bool),
		mqttClient:     messageHandler,
		commandService: cmdService,
		tokenService:   tokenService,
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
		payload := commands[i].Payload
		if commands[i].Payload != "" && commands[i].TokenRequired {
			newPayload, err := attachTokenToCmd(s.ctx, s.tokenService, commands[i].Payload)
			if err != nil {
				fmt.Println("could not parse payload for cmd ", commands[i].Id)
				_ = s.commandService.UpdateStatus(
					s.ctx,
					commands[i].Id.Hex(),
					command.ERRORED,
					err.Error(),
					"",
				)

				continue
			}

			payload = *newPayload
		}

		fmt.Printf("sending command with payload %s to device %s\n", payload, commands[i].DeviceName)
		topic := commands[i].DeviceName + "/cmd"
		encoded, err := encodeCommandMessage(
			commands[i].Id.Hex(),
			commands[i].CommandActionName,
			payload,
			commands[i].Priority,
			commands[i].RequiredOnline,
		)
		if err != nil {
			fmt.Println("could not encode command: ", err.Error())
			continue
		}

		err = s.mqttClient.Publish(encoded, topic, 1)
		if err != nil {
			_ = s.commandService.UpdateStatus(
				s.ctx,
				commands[i].Id.Hex(),
				command.ERRORED,
				fmt.Sprintf(
					"could not send to device because: %s",
					err.Error(),
				),
				"",
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

	fmt.Println("received ack from: ", splittedTopic[0])

	id, state, errorMsg, data, err := decodeAckMessage(string(msg.Payload))
	if err != nil {
		fmt.Println("could not decode ackMessage: ", err.Error())
		return
	}

	err = s.commandService.UpdateStatus(
		s.ctx,
		*id,
		*state,
		*errorMsg,
		*data,
	)

	if err != nil {
		fmt.Println("could not update command status: ", err.Error())
		return
	}
}

func encodeCommandMessage(
	id string,
	messageAction string,
	payload string,
	priority uint,
	reqiuredOnline bool,
) (string, error) {
	type packed struct {
		Id             string
		MessageAction  string
		Payload        string
		Priority       uint
		ReqiuredOnline bool
	}

	b, err := msgpack.Marshal(&packed{
		Id:             id,
		MessageAction:  messageAction,
		Payload:        payload,
		Priority:       priority,
		ReqiuredOnline: reqiuredOnline,
	})
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func decodeAckMessage(
	hexData string,
) (*string, *command.CommandStatus, *string, *string, error) {
	type unpacked struct {
		Id       string
		State    string
		ErrorMsg string
		Data     string
	}

	data, err := hex.DecodeString(hexData)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("could not decode data '%s' because %s",
			hexData,
			err.Error())
	}

	var msg unpacked
	err = msgpack.Unmarshal(data, &msg)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("could not unmarshal message: %s", err.Error())
	}

	state := command.StatusFromString(msg.State)
	if state == nil {
		return nil, nil, nil, nil, fmt.Errorf("state not found: %s", msg.State)
	}

	return &msg.Id, state, &msg.ErrorMsg, &msg.Data, nil
}

func attachTokenToCmd(ctx context.Context, tokenService *token.Service, payload string) (*string, error) {
	var cmdPayload map[string]any
	err := json.Unmarshal([]byte(payload), &cmdPayload)
	if err != nil {
		return nil, err
	}

	otp, err := tokenService.GenerateNewToken(ctx)
	if err != nil {
		return nil, err
	}

	cmdPayload["otp"] = otp.Token

	newPayload := createKeyValuePairs(cmdPayload)
	return &newPayload, nil
}

func createKeyValuePairs(m map[string]any) string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "{")
	for key, value := range m {
		fmt.Fprintf(b, "\"%s\":\"%s\",", key, value)
	}

	s := b.String()
	s = s[:len(s)-1]
	s = s + "}"

	return s
}
