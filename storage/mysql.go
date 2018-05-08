// mqttstore.storage
package storage

import (
	"log"

	"github.com/eclipse/paho.mqtt.golang"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

const insertQuery = "INSERT INTO message (topic, payload) VALUES (?, ?)"

// Storage object with Connection reference and worker Channel
type Storage struct {
	Connection *sqlx.DB
	Channel		chan mqtt.Message
}

// Initialize the connection with the database, return a Storage struct or error
func Initialize(engine string, dsn string) (storage Storage, e error) {
	db, e := sqlx.Connect(engine, dsn)
	if e != nil {
		return Storage{Connection: nil}, e
	}
	log.Println("db.connected:", dsn)
	s := Storage{Connection: db}

	s.Channel = make(chan mqtt.Message)
	go s.Worker()

	return s, nil
}

// Save an mqtt.Message into database
func (storage Storage) save(msg mqtt.Message) error {
	tx := storage.Connection.MustBegin()
	tx.Exec(insertQuery, msg.Topic(), msg.Payload())
	return tx.Commit()
}

// Wait for messages in storage Channel and save them in database
func (storage Storage) Worker() {
	for {
		msg := <-storage.Channel
		storage.save(msg)
	}
}
