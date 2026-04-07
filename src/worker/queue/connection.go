package queue

import (
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

var conn *amqp.Connection
var ch *amqp.Channel

func InitBroker() error {
	var err error

	for i := 1; i <= 10; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@broker:5672/")
		if err == nil {
			break
		}
		log.Printf("Tentativa %d/10: broker indisponível, aguardando 5s...\n", i)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	if err != nil {
		return err
	}

	log.Println("Conectado ao broker com sucesso")
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