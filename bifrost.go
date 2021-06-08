package bifrost

import (
	"fmt"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

/**
* - MQTT hostname
* - MQTT port
* - MQTT username
* - MQTT password
* - MQTT topic
*
* - HTTP host
* - HTTP port
* - HTTP requestType
 */
type Bridge interface {
	Connect(topic string)
	Disconnect()
}

type BridgeOptions struct {
	MQTTHost     string
	MQTTPort     int16
	MQTTUsername string
	MQTTPassword string
	HTTPHost     string
	HTTPMethod   string
}

type bridge struct {
	options    BridgeOptions
	mqttClient mqtt.Client
}

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

func NewBridgeOptions() *BridgeOptions {
	o := &BridgeOptions{
		MQTTHost:     "127.0.0.1",
		MQTTPort:     1883,
		MQTTUsername: "",
		MQTTPassword: "",
		HTTPHost:     "",
		HTTPMethod:   "POST",
	}

	return o
}

func (o *BridgeOptions) AddMQTTHost(host string) *BridgeOptions {
	o.MQTTHost = host

	return o
}

func (o *BridgeOptions) AddMQTTUser(username string, password string) *BridgeOptions {
	o.MQTTUsername = username
	o.MQTTPassword = password

	return o
}

func (o *BridgeOptions) AddHTTPHost(host string) *BridgeOptions {
	o.HTTPHost = host

	return o
}

func NewBridge(o *BridgeOptions) Bridge {
	b := &bridge{}
	b.options = *o

	mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%d", o.MQTTHost, o.MQTTPort)).SetClientID("bifrost_client")
	// mqttOpts.KeepAlive(60 * time.Second)

	mqttOpts.SetDefaultPublishHandler(f)
	mqttOpts.SetPingTimeout(1 * time.Second)

	b.mqttClient = mqtt.NewClient(mqttOpts)

	return b
}

func (b *bridge) Connect(topic string) {
	if token := b.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	if token := b.mqttClient.Subscribe("test-topic", 0, func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("* [%s] %s\n", msg.Topic(), string(msg.Payload()))
	}); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
}

func (b *bridge) Disconnect() {
	if token := b.mqttClient.Unsubscribe("testtopic/#"); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	b.mqttClient.Disconnect(2000)
}
