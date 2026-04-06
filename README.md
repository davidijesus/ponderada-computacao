# sensor-pipeline

Sistema de ingestão de eventos de sensores com arquitetura desacoplada utilizando **Go, RabbitMQ, PostgreSQL, Docker e k6**.

## Visão geral

O sistema recebe eventos de dispositivos e os processa de forma assíncrona.

Fluxo da aplicação:

```
Device → API → RabbitMQ → Worker → PostgreSQL
```

## Estrutura de pastas

```
src/
  api/               # API HTTP (producer)
    handlers/
    models/
    queue/
    router/
    tests/

  worker/            # Worker (consumer da fila)
    queue/
    storage/
    models/
    tests/

  db/
    schema.sql       # criação da tabela

  tests/
    loadtest/        # scripts k6
```

## Como rodar o projeto

### 1. Subir todos os serviços

```bash
docker compose up --build
```

### 2. Enviar requisição

em outro terminal:

```bash
curl -X POST http://localhost:8080/events \
  -H "Content-Type: application/json" \
  -d '{
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
  }'
```

### 3. Verificar banco

```bash
docker exec -it database psql -U pipeline_user -d pipeline_db
```

```sql
SELECT * FROM sensor_events;
```

Para sair: `\q`

## Testes unitários

testes da API:
```bash
cd src/api
go test ./...
```

testes do worker:
```bash
cd src/worker
go test ./...
```

## Teste de carga (k6)

Com os serviços no ar (`docker compose up --build`):

### Cenários disponíveis

* `sensor_load_test.js` → carga progressiva forte (até ~28 VUs)
* `sensor_stress_test.js` → stress alto calibrado para aprovação (até ~55 VUs)
* `sensor_pico_test.js` → picos agressivos e curtos (até ~160 VUs)
* `sensor_resistencia_test.js` → carga sustentada forte em longa duração (até ~45 VUs)
* `sensor_breakpoint_test.js` → rampa forte de limite de capacidade

### Como rodar

Padrão (roda `sensor_load_test.js`):

```bash
docker compose --profile loadtest run --rm k6
```

Escolher cenário específico:

```bash
K6_SCRIPT=sensor_load_test.js docker compose --profile loadtest run --rm k6
K6_SCRIPT=sensor_stress_test.js docker compose --profile loadtest run --rm k6
K6_SCRIPT=sensor_pico_test.js docker compose --profile loadtest run --rm k6
K6_SCRIPT=sensor_resistencia_test.js docker compose --profile loadtest run --rm k6
```

Rodar todos os cenários de uma vez (sequencial):

```bash
bash run_all_loadtests.sh
```

### Onde os resultados são salvos

Cada execução salva automaticamente dois arquivos em `src/tests/loadtest/results`:

* `*_summary.json` (resumo final do teste)
* `*_metrics.json` (métricas detalhadas)

## Arquitetura

API:
* recebe HTTP
* valida dados
* envia para fila

RabbitMQ:
* fila de mensagens
* desacoplamento

Worker:
* lê da fila
* converte JSON
* salva no banco

PostgreSQL:
* persistência
