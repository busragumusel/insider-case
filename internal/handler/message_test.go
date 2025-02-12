package handler

import (
	"encoding/json"
	"github.com/busragumusel/insider-case/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockMessageService struct{}

func (m *mockMessageService) StartProcess() {}

func (m *mockMessageService) StopProcess() {}

func (m *mockMessageService) Retrieve() ([]entity.Message, error) {
	return []entity.Message{
		{
			ID:          1,
			PhoneNumber: "+90532434532",
		},
	}, nil
}

func TestStartProcess(t *testing.T) {
	service := &mockMessageService{}
	handler := NewMessageHandler(service)

	req, err := http.NewRequest("GET", "/start", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.StartProcess(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestStopProcess(t *testing.T) {
	service := &mockMessageService{}
	handler := NewMessageHandler(service)

	req, err := http.NewRequest("GET", "/stop", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.StopProcess(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestRetrieve(t *testing.T) {
	service := &mockMessageService{}
	handler := NewMessageHandler(service)

	req, err := http.NewRequest("GET", "/messages", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()
	handler.Retrieve(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "data")
}
