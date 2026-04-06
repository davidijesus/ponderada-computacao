package queue

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitBroker() error {
	var err error

	conn, err = amqp.Dial("amqp://guest:guest@broker:5672/")
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	return nil
}

func GetChannel() *amqp.Channel {
	return ch
}

func CloseBroker() {
	if ch != nil {
		ch.Close()
	}
	if conn != nil {
		conn.Close()
	}
}
