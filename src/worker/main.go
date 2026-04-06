package main

import (
	"log"
	"worker/queue"
	"worker/storage"
)

func main() {
	if err := queue.InitBroker(); err != nil {
		log.Fatal("Falha ao conectar no broker:", err)
	}
	defer queue.CloseBroker()

	if err := storage.InitDB(); err != nil {
		log.Fatal("Falha ao conectar no banco:", err)
	}

	if err := queue.StartWorker("sensor_events"); err != nil {
		log.Fatal("Falha ao iniciar worker:", err)
	}
}
