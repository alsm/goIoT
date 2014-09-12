package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	MQTT "git.eclipse.org/gitroot/paho/org.eclipse.paho.mqtt.golang.git"
	"github.com/kidoman/embd"
	_ "github.com/kidoman/embd/host/all"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Data struct {
	Load1   float64 `json:"cpuload"`
	Load5   float64 `json:"cpuload5"`
	Load15  float64 `json:"cpuload15"`
	CPUTemp float64 `json:"cputemp"`
}

type Payload struct {
	DataPoints Data      `json:"d"`
	Timestamp  time.Time `json:"ts"`
}

var (
	host              Platform
	configFile        *string
	config            *Config
	quickstartBaseURL string = "http://quickstart.internetofthings.ibmcloud.com/#/device/"
)

func init() {
	configFile = flag.String("conf", "", "IoT app configuration file")
	flag.Parse()
	switch h, _, err := embd.DetectHost(); h {
	case embd.HostRPi:
		host = NewRPi()
	case embd.HostBBB:
		host = NewBeagleBone()
	default:
		if err != nil {
			panic(err)
		}
	}
}

func actionHandler(client *MQTT.MqttClient, message MQTT.Message) {
	fmt.Println("Received action message on", message.Topic(), "-", string(message.Payload()))
	action := strings.ToLower(string(message.Payload()))
	switch action {
	case "off":
		host.LedsOff()
	case "on":
		host.LedsOn()
	case "toggle":
		host.LedsToggle()
	case "slide":
		host.LedsCycle(3)
	}
}

func SendData(client MQTT.Client, endChan chan struct{}) {
	ticker := time.NewTicker(1 * time.Second)
	for {
		select {
		case <-endChan:
			return
		case <-ticker.C:
			var p Payload
			p.DataPoints.Load1, p.DataPoints.Load5, p.DataPoints.Load15 = getLoadAvg()
			p.DataPoints.CPUTemp = getCPUTemp()
			p.Timestamp = time.Now()
			payloadBytes, err := json.Marshal(p)
			if err == nil {
				client.Publish(config.PubTopic, 0, false, payloadBytes)
			} else {
				fmt.Println(err.Error())
			}
		}
	}
}

func main() {
	var err error
	endChan := make(chan struct{})
	host.LedsOff()

	config = ParseConfig(*configFile)

	fmt.Println("Device ID:", config.DeviceID)
	fmt.Println("Connecting to MQTT broker:", config.BrokerAddress)

	opts := MQTT.NewClientOptions().AddBroker(config.BrokerAddress).SetClientId(config.ClientID)
	if !config.QuickStart {
		tlsConfig := &tls.Config{InsecureSkipVerify: true, ClientAuth: tls.NoClientCert}
		opts.SetUsername(config.Username).SetPassword(config.AuthToken).SetTlsConfig(tlsConfig)
	}
	client := MQTT.NewClient(opts)
	_, err = client.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	fmt.Println("Connected")
	host.LedsCycle(3)

	go SendData(client, endChan)
	fmt.Println("Sending Data")

	if config.QuickStart {
		fmt.Println("Go to the following link to see your device data;")
		fmt.Println(quickstartBaseURL + config.DeviceID + "/sensor/")
	} else {
		var token *MQTT.SubscribeToken
		fmt.Println("Subscribing for action messages")
		token, err = client.Subscribe("iot-2/cmd/+/fmt/text", 0, actionHandler)
		if err != nil {
			fmt.Println("Error subscribing for action messages")
		}
		token.Wait()
		for topic, qos := range token.Results() {
			fmt.Println("Subscribed to", topic, "at Qos", qos)
		}
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	host.Close()
}
