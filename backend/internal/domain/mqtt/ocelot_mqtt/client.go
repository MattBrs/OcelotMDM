package ocelot_mqtt

import (
	"fmt"
	"strings"

	"github.com/MattBrs/OcelotMDM/internal/domain/mqtt/paho_mqtt"
)

type MqttClient struct {
	pahoClient    paho_mqtt.MqttClient
	AckMessages   chan ChanMessage
	LogMessages   chan ChanMessage
	OtherMessages chan ChanMessage
	doneChannel   chan bool
}

type ChanMessage struct {
	Topic   string
	Payload []byte
}

func NewMqttClient(server string, port uint) *MqttClient {
	client := MqttClient{
		pahoClient:    paho_mqtt.NewMqttClient(server, port),
		AckMessages:   make(chan ChanMessage, 1000),
		LogMessages:   make(chan ChanMessage, 1000),
		OtherMessages: make(chan ChanMessage, 1000),
		doneChannel:   make(chan bool),
	}

	go func() {
		demuxChannels(&client)
	}()

	return &client
}

func demuxChannels(client *MqttClient) {
	for {
		select {
		case <-client.doneChannel:
			return
		case msg := <-client.pahoClient.Messages:
			fwMsg := ChanMessage{
				Topic:   msg.Topic(),
				Payload: make([]byte, len(msg.Payload())),
			}
			copy(fwMsg.Payload, msg.Payload())

			if strings.Contains(msg.Topic(), "ack") {
				tryEnqueue(client.AckMessages, fwMsg, "ack")
			} else if strings.Contains(msg.Topic(), "logs") {
				tryEnqueue(client.LogMessages, fwMsg, "logs")
			} else {
				tryEnqueue(client.OtherMessages, fwMsg, "other")
			}
		}
	}
}

func tryEnqueue(ch chan ChanMessage, msg ChanMessage, name string) {
	select {
	case ch <- msg:
	default:
		fmt.Println("dropped msg for channel: ", name)
	}
}

func (client *MqttClient) Connect() error {
	return client.pahoClient.Connect()
}

func (client *MqttClient) Close() {
	client.doneChannel <- true

	client.pahoClient.Close()

	close(client.doneChannel)
	close(client.AckMessages)
	close(client.LogMessages)
	close(client.OtherMessages)
}

func (client *MqttClient) Subscribe(topic string, qos byte) error {
	return client.pahoClient.Subscribe(topic, qos)
}

func (client *MqttClient) Publish(
	message string,
	topic string,
	qos byte,
) error {
	return client.pahoClient.Publish(message, topic, qos)
}
