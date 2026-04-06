package queue

import (
	"errors"
	"testing"
	"time"
	"worker/queue"
	"worker/models"
)

func validEventPayload() []byte {
	return []byte(`{
		"device_id": 1,
		"timestamp": "2026-03-20T13:15:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 25.3
		}
	}`)
}

func invalidEventPayload() []byte {
	return []byte(`{
		"device_id":
	}`)
}

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("esperava nil, recebeu erro: %v", err)
	}
}

func assertError(t *testing.T, err error) {
	t.Helper()
	if err == nil {
		t.Fatal("esperava erro, recebeu nil")
	}
}

func assertEventFields(t *testing.T, got models.SensorEvent) {
	t.Helper()

	if got.DeviceID != 1 {
		t.Errorf("device_id incorreto: esperado=%d, recebido=%d", 1, got.DeviceID)
	}

	expectedTime := time.Date(2026, 3, 20, 13, 15, 0, 0, time.UTC)
	if !got.Timestamp.Equal(expectedTime) {
		t.Errorf("timestamp incorreto: esperado=%v, recebido=%v", expectedTime, got.Timestamp)
	}

	if got.Sensor == nil {
		t.Fatal("sensor não deveria ser nil")
	}

	if got.Sensor.Kind != "temperature" {
		t.Errorf("sensor.type incorreto: esperado=%v, recebido=%v", "temperature", got.Sensor.Kind)
	}

	if got.Sensor.Unit != "celsius" {
		t.Errorf("sensor.unit incorreto: esperado=%v, recebido=%v", "celsius", got.Sensor.Unit)
	}

	if got.Reading == nil {
		t.Fatal("reading não deveria ser nil")
	}

	if got.Reading.Kind != "analog" {
		t.Errorf("reading.value_type incorreto: esperado=%v, recebido=%v", "analog", got.Reading.Kind)
	}

	if got.Reading.Value != 25.3 {
		t.Errorf("reading.value incorreto: esperado=%v, recebido=%v", 25.3, got.Reading.Value)
	}
}

func TestProcessEvent_Success(t *testing.T) {
	payload := validEventPayload()

	var saved models.SensorEvent
	saveFn := func(e models.SensorEvent) error {
		saved = e
		return nil
	}

	err := queue.ProcessEvent(payload, saveFn)
	assertNoError(t, err)
	assertEventFields(t, saved)
}

func TestProcessEvent_InvalidJSON(t *testing.T) {
	payload := invalidEventPayload()

	called := false
	saveFn := func(e models.SensorEvent) error {
		called = true
		return nil
	}

	err := queue.ProcessEvent(payload, saveFn)
	assertError(t, err)

	if called {
		t.Fatal("saveFn não deveria ter sido chamada para JSON inválido")
	}
}

func TestProcessEvent_SaveError(t *testing.T) {
	payload := validEventPayload()

	expectedErr := errors.New("erro ao persistir no banco")
	saveFn := func(e models.SensorEvent) error {
		return expectedErr
	}

	err := queue.ProcessEvent(payload, saveFn)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("erro incorreto: esperado=%v, recebido=%v", expectedErr, err)
	}
}
