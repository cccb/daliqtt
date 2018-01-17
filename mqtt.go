package main

import (
	"fmt"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Host     string
	User     string
	Password string

	BaseTopic string
}

func (config MqttConfig) BrokerUri() string {
	uri := "tcp://"
	if config.User != "" {
		uri += config.User

		if config.Password != "" {
			uri += ":" + config.Password
		}

		uri += "@"
	}

	uri += config.Host

	return uri
}

func DialMqtt(config *MqttConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(config.BrokerUri()).
		SetClientID("daliqtt")

	opts.SetPingTimeout(1 * time.Second)
	opts.SetKeepAlive(2 * time.Second)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("TOPIC: %s\n", msg.Topic()) // String
		fmt.Printf("MSG: %s\n", msg.Payload()) // []byte

		if msg.Topic() == "fnord/FOO" {
			client.Publish("fnord/FOO_BAM", 0, false, "foooooo")
		}
	})

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	// Subscribe to topic
	topic := config.BaseTopic + "/#"
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	return client, nil
}
