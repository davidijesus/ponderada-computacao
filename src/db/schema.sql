CREATE TABLE IF NOT EXISTS sensor_events (
    id BIGSERIAL PRIMARY KEY,

    device_id INTEGER NOT NULL,

    timestamp TIMESTAMPTZ NOT NULL,

    sensor_type TEXT NOT NULL,
    sensor_unit TEXT NOT NULL,

    reading_type TEXT NOT NULL,
    value DOUBLE PRECISION NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
