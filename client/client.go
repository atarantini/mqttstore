// mqttstore.client
package client

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eclipse/paho.mqtt.golang"

	"github.com/atarantini/mqttstore/storage"
)

var store *storage.Storage

var statsIntervalSeconds = 10
var statsMessagesTotal = 0

// Start MQTT client by connecting to server, subscribe into a topic and save messages
func Start(mqttServer string, mqttServerPort string, mqttTopic string, mqttClientID string, storage *storage.Storage) {
	store = storage

	// Set MQTT client options
	opts := mqtt.NewClientOptions().AddBroker("tcp://" + fmt.Sprintf("%s:%s", mqttServer, mqttServerPort))
	opts.SetClientID(mqttClientID)
	opts.SetAutoReconnect(true)
	opts.SetConnectionLostHandler(connLostHandler)
	opts.SetDefaultPublishHandler(mqttMessageHandler) // This handler will send messages to store.Channel
	opts.OnConnect = func(c mqtt.Client) {
		// Connection OK, subscribe to topic
		log.Println("mqtt.topic:", mqttTopic)
		c.Subscribe(mqttTopic, 0, nil)
		log.Println("mqtt.client: ignoring_retained_messages (this can take a long time depending on your broker stored retained messages)")
	}

	//
	// Create client and connect to server
	//
	c := mqtt.NewClient(opts)
	log.Printf("mqtt.server: %s:%s (%s)", mqttServer, mqttServerPort, mqttClientID)

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalln("mqtt.error:", token.Error())
		os.Exit(100)
	}

	//
	// Stats every `statsIntervalSeconds`
	//
	ticker := time.NewTicker(time.Second * time.Duration(statsIntervalSeconds))
	go func() {
		for t := range ticker.C {
			printStats(t)
		}
	}()

}

// Handle each MQTT message
var mqttMessageHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	// Ignore retained messages
	if msg.Retained() {
		return
	}
	// Send message to store.Channel
	store.Channel <- msg

	// Increment stats total counter
	statsMessagesTotal++
}

// Connection lost, print error
func connLostHandler(c mqtt.Client, err error) {
	log.Println(fmt.Sprintf("mqtt.error.connection_lost: %v", err))
}

// Print stats
func printStats(t time.Time) {
	if statsMessagesTotal == 0 {
		return
	}
	log.Println("mqtt.stats.received_messages_total:", statsMessagesTotal)
}