package main

import (
	"log"
	"api/queue"
	"api/router"
)

func main() {
	err := queue.InitBroker()
	if err != nil {
		log.Fatal("Falha ao conectar no broker:", err)
	}

	r := router.SetupRouter()
	r.Run(":8080")
}
