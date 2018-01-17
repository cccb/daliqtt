package main

import (
	"flag"
	"log"
	"os"

	"github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	Mqtt         *MqttConfig
	LichtCgiBase string
}

func parseFlags() *Config {
	host := flag.String("host", "localhost:1883", "MQTT broker host")
	user := flag.String("user", "", "MQTT broker host")
	password := flag.String("password", "", "MQTT broker host")
	baseTopic := flag.String("topic", "dali", "MQTT base topic")

	lichtCgiBase := flag.String("lichtcgi", "http://dali", "licht.cgi server")

	flag.Parse()

	mqttConfig := &MqttConfig{
		Host:     *host,
		User:     *user,
		Password: *password,

		BaseTopic: *baseTopic,
	}

	config := &Config{
		Mqtt:         mqttConfig,
		LichtCgiBase: *lichtCgiBase,
	}

	return config
}

func main() {
	ctrl := make(chan os.Signal, 1)

	log.Println("Starting dali to mqtt bridge")

	// Initialize configuration
	config := parseFlags()

	// Initialize MQTT connection
	mqtt.ERROR = log.New(os.Stdout, "", 0)
	actions, dispatch, err := DialMqtt(config.Mqtt)
	if err != nil {
		panic(err)
	}

	// So far so good. Let now the lights service
	// take over and handle requests.
	svc := NewLightsSvc(config.LichtCgiBase)
	go svc.Handle(actions, dispatch)

	<-ctrl
}
