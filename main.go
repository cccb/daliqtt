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

func main() {

	fmt.Println("Starting dali to mqtt bridge")

	// Initialize configuration
	config := parseFlags()

	lights := NewLightsState()
	lights.Refresh()

	fmt.Println(config)

}
