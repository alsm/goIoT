package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Config struct {
	QuickStart    bool
	BrokerAddress string
	ClientID      string
	Username      string
	Organisation  string `json:"org"`
	Type          string `json:"type"`
	DeviceID      string `json:"device-id"`
	AuthToken     string `json:"auth-token"`
	DeviceName    string `json:"device-name"`
	PubTopic      string `json:"publish-topic"`
}

var (
	interfaces []net.Interface
)

func init() {
	interfaces, _ = net.Interfaces()
}

func ParseConfig(confFile string) *Config {
	var config Config

	if confFile == "" {
		return DefaultConfig(&config)
	}

	file, err := os.Open(confFile)
	if err != nil {
		return DefaultConfig(&config)
	}
	decoder := json.NewDecoder(file)

	err = decoder.Decode(&config)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Using quickstart configuration")
		return DefaultConfig(&config)
	}

	if config.Organisation == "" || config.AuthToken == "" {
		fmt.Println("Organisation or AuthToken is not set. Using quickstart configuration")
		return DefaultConfig(&config)
	}

	if config.DeviceID == "" {
		hwAddr := interfaces[1].HardwareAddr.String()
		config.DeviceID = strings.Replace(hwAddr, ":", "", -1)
	}

	if config.Type == "" {
		config.Type = "iotsample-raspberrypi"
	}

	if config.DeviceName == "" {
		config.DeviceName = "GoIoT"
	}

	if config.PubTopic == "" {
		config.PubTopic = "iot-2/evt/status/fmt/json"
	}

	config.Username = "use-token-auth"
	config.BrokerAddress = "ssl://" + config.Organisation + ".messaging.internetofthings.ibmcloud.com:8883"
	config.ClientID = strings.Join([]string{"d", config.Organisation, config.Type, config.DeviceID}, ":")
	return &config
}

func DefaultConfig(confVar *Config) *Config {
	hwAddr := interfaces[1].HardwareAddr.String()
	confVar.QuickStart = true
	confVar.DeviceID = strings.Replace(hwAddr, ":", "", -1)
	confVar.Type = "iotsample-raspberrypi"
	confVar.PubTopic = "iot-2/evt/status/fmt/json"
	confVar.BrokerAddress = "tcp://messaging.quickstart.internetofthings.ibmcloud.com:1883"
	confVar.ClientID = strings.Join([]string{"d", "quickstart", confVar.Type, confVar.DeviceID}, ":")
	confVar.DeviceName = "GoIoT"
	return confVar
}
