package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

/*
 Decode an incoming mqtt message and create an
 action from it's topic and payload
*/
func decodeMessage(msg mqtt.Message) (Action, error) {
	// Decode topic
	tokens := strings.Split(msg.Topic(), "/")
	actionType := tokens[len(tokens)-1]

	// Decode payload
	var payload interface{}
	var err error
	switch actionType {
	case SET_LIGHT_VALUE_REQUEST:
		var lightValue LightValuePayload
		err = json.Unmarshal(msg.Payload(), &lightValue)
		payload = lightValue
	}

	// Make action
	action := Action{
		Type:    actionType,
		Payload: payload,
	}

	return action, err
}

func DialMqtt(config *MqttConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(config.BrokerUri()).
		SetClientID("daliqtt")

	opts.SetPingTimeout(1 * time.Second)
	opts.SetKeepAlive(2 * time.Second)

	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		action, err := decodeMessage(msg)
		if err != nil {
			log.Println("Error while decoding message:", err)
			return
		}

		fmt.Println("Incoming action:", action)

		if action.Type == SET_LIGHT_VALUE_REQUEST {
			request := action.Payload.(LightValuePayload)
			fmt.Println("Setting light", request.Id, "to", request.Value)
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
