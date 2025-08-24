package paho_mqtt

import (
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client mqtt.Client
}

var messageHandler = func(client mqtt.Client, message mqtt.Message) {
	fmt.Printf(
		"received message %s from topic %s\n",
		message.Payload(),
		message.Topic(),
	)
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("connected to mqtt broker ", time.Now().UnixMilli())
}

var connectionLostHandler mqtt.ConnectionLostHandler = func(
	client mqtt.Client,
	err error,
) {
	fmt.Println("connection lost. cause: ", err.Error())
}

func NewMqttClient(server string, port uint) MqttClient {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", server, port))
	opts.SetClientID("ocelot_backend")
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectionLostHandler
	opts.SetDefaultPublishHandler(messageHandler)

	return MqttClient{
		client: mqtt.NewClient(opts),
	}
}

func (mc MqttClient) Publish(message string, topic string, qos byte) {
	mc.client.Publish(topic, qos, true, message)
}

func (mc MqttClient) Subscribe(topic string, qos byte) {
	subToken := mc.client.Subscribe(topic, qos, nil)
	subToken.Wait()
	fmt.Println("subscribed successfully to topic: ", topic)
}

func (mc MqttClient) Connect() error {
	connectToken := mc.client.Connect()
	connectToken.Wait()

	return connectToken.Error()
}
