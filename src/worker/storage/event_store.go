package storage

import (
	"worker/models"
)

func InsertEvent(e models.SensorEvent) error {
	query := `
		INSERT INTO sensor_events (
			device_id,
			timestamp,
			sensor_type,
			sensor_unit,
			reading_type,
			value
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := DB.Exec(
		query,
		e.DeviceID,
		e.Timestamp,
		e.Sensor.Kind,
		e.Sensor.Unit,
		e.Reading.Kind,
		e.Reading.Value,
	)

	return err
}
