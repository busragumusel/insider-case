package handler

import (
	"context"
	"encoding/json"
	"github.com/busragumusel/insider-case/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockMessageService struct{}

func (m *mockMessageService) StartProcess(_ context.Context) {}

func (m *mockMessageService) StopProcess(_ context.Context) {}

func (m *mockMessageService) Retrieve(_ context.Context, status string) ([]entity.Message, error) {
	return []entity.Message{
		{
			ID:          1,
			PhoneNumber: "+90532434532",
			Content:     "Test Message",
			Status:      status,
			CreatedAt:   time.Now().Add(-10 * time.Minute),
			SentAt:      time.Time{}, // Default zero time (not sent yet)
		},
		{
			ID:          2,
			PhoneNumber: "+905551111111",
			Content:     "Another Test",
			Status:      status,
			CreatedAt:   time.Now().Add(-5 * time.Minute),
			SentAt:      time.Now(),
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
