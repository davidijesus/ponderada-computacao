package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gin-gonic/gin"
	"api/router"
	"api/queue"
)

func init() {
	queue.Publish = func(channelName string, payload interface{}) error {
		return nil
	}
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.SetupRouter()
}

func doRequest(r *gin.Engine, method, path, body string) (*httptest.ResponseRecorder, map[string]interface{}) {
	req, _ := http.NewRequest(method, path, bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)

	return w, resp
}

func checkError(t *testing.T, w *httptest.ResponseRecorder, resp map[string]interface{}, wantStatus int, wantMsg string) {
	if w.Code != wantStatus {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", wantStatus, w.Code)
	}
	if resp["error"] != wantMsg {
		t.Errorf("mensagem incorreta: esperado=%v, recebido=%v", wantMsg, resp["error"])
	}
}

func checkSuccess(t *testing.T, w *httptest.ResponseRecorder, resp map[string]interface{}) {
	if w.Code != http.StatusOK {
		t.Errorf("status incorreto: esperado=%d, recebido=%d", http.StatusOK, w.Code)
	}
	expected := "Evento recebido com sucesso"
	if resp["message"] != expected {
		t.Errorf("mensagem incorreta: esperado=%v, recebido=%v", expected, resp["message"])
	}
	if resp["data"] == nil {
		t.Errorf("campo data não encontrado")
	}
}

func TestHandleSensorEvent_Success(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkSuccess(t, w, resp)
}

func TestHandleSensorEvent_InvalidJSON(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id":
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "Payload inválido")
}

func TestHandleSensorEvent_InvalidDeviceID(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 0,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "device_id é obrigatório")
}

func TestHandleSensorEvent_EmptySensorType(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "sensor.type é obrigatório")
}

func TestHandleSensorEvent_EmptySensorUnit(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": ""
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "sensor.unit é obrigatório")
}

func TestHandleSensorEvent_EmptyValueType(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "reading.value_type é obrigatório")
}

func TestHandleSensorEvent_InvalidTimestampFormat(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "data-errada",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "Payload inválido")
}

func TestHandleSensorEvent_InvalidTimestampType(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": 123456,
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "Payload inválido")
}

func TestHandleSensorEvent_MissingTimestamp(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "timestamp é obrigatório")
}

func TestHandleSensorEvent_EmptyTimestamp(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "Payload inválido")
}

func TestHandleSensorEvent_InvalidValueType(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		},
		"reading": {
			"value_type": "analogo",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "value_type deve ser 'analog' ou 'discrete'")
}

func TestHandleSensorEvent_MissingSensor(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"reading": {
			"value_type": "analog",
			"value": 23.7
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "sensor é obrigatório")
}

func TestHandleSensorEvent_MissingReading(t *testing.T) {
	r := setupTestRouter()

	payload := `{
		"device_id": 1,
		"timestamp": "2026-03-17T14:30:00Z",
		"sensor": {
			"type": "temperature",
			"unit": "celsius"
		}
	}`

	w, resp := doRequest(r, "POST", "/events", payload)
	checkError(t, w, resp, http.StatusBadRequest, "reading é obrigatório")
}
