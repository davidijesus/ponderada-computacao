package models

import "time"

type SensorMeta struct {
	Kind string `json:"type"`
	Unit string `json:"unit"`
}

type SensorData struct {
	Kind  string  `json:"value_type"`
	Value float64 `json:"value"`
}

type SensorEvent struct {
	DeviceID  int         `json:"device_id"`
	Timestamp time.Time   `json:"timestamp"`
	Sensor    *SensorMeta `json:"sensor"`
	Reading   *SensorData `json:"reading"`
}
