package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection

func InitBroker() error {
	var err error

	conn, err = amqp.Dial("amqp://guest:guest@broker:5672/")
	if err != nil {
		return err
	}

	return nil
}

func GetConnection() *amqp.Connection {
	return conn
}

func CloseBroker() {
	if conn != nil {
		conn.Close()
	}
}
