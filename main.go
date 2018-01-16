package main

import (
	"flag"
	"fmt"
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

func main() {

	fmt.Println("Starting dali to mqtt bridge")

	// Initialize configuration
	config := parseFlags()

	fmt.Println(config)

}
