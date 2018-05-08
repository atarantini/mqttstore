// mqttstore
//
// Store MQTT messages into MySQL database
package main

import (
	"flag"
	"log"

	"github.com/satori/go.uuid"

	"github.com/atarantini/mqttstore/client"
	"github.com/atarantini/mqttstore/storage"
)

var dsn string
var mqttServer string
var mqttServerPort string
var mqttTopic string
var mqttClientID string

// Command line arguments constants
const claDefaultDsn = "test:test@/mqtt"

func main() {
	// Parse command line arguments
	flag.StringVar(&dsn, "dsn", claDefaultDsn, "Data source name (e.g.: user:password@hostname/database_name)")
	flag.StringVar(&mqttServer, "host", "iot.eclipse.org", "MQTT server hostname or IP address")
	flag.StringVar(&mqttServerPort, "port", "1883", "MQTT server port")
	flag.StringVar(&mqttTopic, "topic", "#", "MQTT topic")
	uuid, _ := uuid.NewV4()
	clientID := uuid.String()
	flag.StringVar(&mqttClientID, "clientid", clientID, "MQTT Client ID")
    flag.Parse()

	// Connect to DB
	store, err := storage.Initialize("mysql", dsn)
	if err != nil {
		log.Println("db.error: Using database '" + dsn + "', set up your database connection setting up a DSN using the '-dsn' parameter:\n\n\t$ mqtt -dsn username:password@localhost/database_name\n")
		log.Fatalln("db.error:", err)
	}

	client.Start(mqttServer, mqttServerPort, mqttTopic, mqttClientID, &store)

	// Block forever (until signal or CTRL+C)
	select {}
}