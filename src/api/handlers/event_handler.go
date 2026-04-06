package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"api/models"
	"api/queue"
)

func requireField(c *gin.Context, value string, field string) bool {
	if value == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": field + " é obrigatório",
		})
		return false
	}
	return true
}

func HandleSensorEvent(c *gin.Context) {
	var event models.SensorEvent

	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Payload inválido",
		})
		return
	}

	if event.DeviceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "device_id é obrigatório",
		})
		return
	}

	if event.Sensor == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "sensor é obrigatório",
		})
		return
	}

	if event.Reading == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "reading é obrigatório",
		})
		return
	}

	if !requireField(c, event.Sensor.Kind, "sensor.type") {
		return
	}

	if !requireField(c, event.Sensor.Unit, "sensor.unit") {
		return
	}

	if !requireField(c, event.Reading.Kind, "reading.value_type") {
		return
	}

	if event.Timestamp.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "timestamp é obrigatório",
		})
		return
	}

	if event.Reading.Kind != "analog" && event.Reading.Kind != "discrete" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "value_type deve ser 'analog' ou 'discrete'",
		})
		return
	}

	if err := queue.Publish("sensor_events", event); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Falha ao enfileirar evento",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Evento recebido com sucesso",
		"data":    event,
	})
}
