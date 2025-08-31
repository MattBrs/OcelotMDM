package paho_mqtt

import (
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client   mqtt.Client
	Messages chan mqtt.Message
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
	opts.CleanSession = false

	return MqttClient{
		client:   mqtt.NewClient(opts),
		Messages: make(chan mqtt.Message, 1000),
	}
}

func (mc MqttClient) Publish(message string, topic string, qos byte) error {
	pubToken := mc.client.Publish(topic, qos, true, message)

	pubToken.Wait()
	return pubToken.Error()
}

func (mc MqttClient) Subscribe(topic string, qos byte) error {
	subToken := mc.client.Subscribe(topic, qos, func(_ mqtt.Client, msg mqtt.Message) {
		select {
		case mc.Messages <- msg:
		default:
			fmt.Println(
				"MQTT queue full, discarded message: ",
				string(msg.Payload()),
			)
		}
	})

	subToken.Wait()
	return subToken.Error()
}

func (mc MqttClient) Connect() error {
	connectToken := mc.client.Connect()
	connectToken.Wait()

	return connectToken.Error()
}

func (mc MqttClient) Close() {
	mc.client.Disconnect(500)
	close(mc.Messages)
}
