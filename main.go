package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Mqtt *MqttConfig
}

func parseFlags() *Config {
	host := flag.String("host", "localhost:1883", "MQTT broker host")
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

func main() {
	ctrl := make(chan os.Signal, 1)

	fmt.Println("Starting dali to mqtt bridge")

	// Initialize configuration
	config := parseFlags()

	lights := NewLightsState()
	lights.Refresh()

	// MQTT test
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	_, err := DialMqtt(config.Mqtt)
	if err != nil {
		panic(err)
	}

	<-ctrl
}
