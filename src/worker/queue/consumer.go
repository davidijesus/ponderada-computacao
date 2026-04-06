package queue

import (
	"log"
	"encoding/json"
	"worker/models"
	"worker/storage"
)

func ProcessEvent(body []byte, save func(models.SensorEvent) error) error {
	var event models.SensorEvent

	if err := json.Unmarshal(body, &event); err != nil {
		return err
	}

	return save(event)
}

func StartWorker(channelName string) error {
	ch := GetChannel()

	q, err := ch.QueueDeclare(
		channelName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Worker ativo, aguardando eventos...")

	for msg := range msgs {
		log.Printf("Evento recebido: %s\n", msg.Body)

		if err := ProcessEvent(msg.Body, storage.InsertEvent); err != nil {
			log.Println("Erro ao processar evento:", err)
			continue
		}

		log.Println("Evento persistido com sucesso")
	}

	return nil
}
