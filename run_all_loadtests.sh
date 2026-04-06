#!/bin/bash

set -e

SCRIPTS=(
  "sensor_load_test.js"
  "sensor_stress_test.js"
  "sensor_pico_test.js"
  "sensor_resistencia_test.js"
  "sensor_breakpoint_test.js"
)

for SCRIPT in "${SCRIPTS[@]}"; do
  echo "Executando: $SCRIPT"
  K6_SCRIPT=$SCRIPT docker compose --profile loadtest run --rm k6
  echo "Concluído: $SCRIPT"
  echo "---"
done
