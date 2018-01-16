package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"
)

type MqttConfig struct {
	Host     string
	User     string
	Password string

	BaseTopic string
}

type Config struct {
	Mqtt *MqttConfig
}

func parseFlags() *Config {
	host := flag.String("host", "localhost", "MQTT broker host")
	user := flag.String("user", "", "MQTT broker host")
	password := flag.String("password", "", "MQTT broker host")
	baseTopic := flag.String("topic", "dali", "MQTT base topic")

	mqttConfig := &MqttConfig{
		Host:     *host,
		User:     *user,
		Password: *password,

		BaseTopic: *baseTopic,
	}

	config := &Config{
		Mqtt: mqttConfig,
	}

	return config
}

// Handle light state
type Light struct {
	Value int
}

type LightsState map[string]Light

/*
 Initialize lights state:
 Assign handles to light value
*/
func NewLightsState() *LightsState {
	state := &LightsState{
		"entry":     Light{0},
		"foh":       Light{0},
		"desk_wall": Light{0},
		"desk_bar":  Light{0},
	}

	return state
}

func (self *LightsState) Set(handle string, value int) {
	(*self)[handle] = Light{value}

	// TODO: Talk to dali
}

func (self *LightsState) Read(handle string) int {
	// TODO: Talk to dali

	return (*self)[handle].Value
}

func (self *LightsState) Refresh() {
	for handle, _ := range *self {
		fmt.Println("Refreshing state of", handle)
		self.Read(handle)
	}
}

func DialMqtt(config *MqttConfig) (mqtt.Client, error) {
	opts := mqtt.NewClientOptions().
		AddBroker("tcp://localhost:1883").
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

	return client, nil
}

func main() {
	ctrl := make(chan os.Signal, 1)

	fmt.Println("Starting dali to mqtt bridge")

	// Initialize configuration
	config := parseFlags()

	lights := NewLightsState()
	lights.Refresh()

	fmt.Println(config)

	// MQTT test
	// mqtt.DEBUG = log.New(os.Stdout, "", 0)
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	client, err := DialMqtt(config.Mqtt)
	if err != nil {
		panic(err)
	}

	// Subscription topic
	topic := config.Mqtt.BaseTopic + "/#"
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	<-ctrl
}
