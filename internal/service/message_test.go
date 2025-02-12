package service

import (
	"encoding/json"
	"errors"
	"github.com/busragumusel/insider-case/internal/model"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepo struct {
	mock.Mock
}

func (m *MockMessageRepo) GetByStatus(status string, limit int) ([]entity.Message, error) {
	args := m.Called(status, limit)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *MockMessageRepo) Update(id uint, status string) error {
	args := m.Called(id, status)
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestRetrieveMessages(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool)

	messages := []entity.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Test Message", Status: "sent"},
	}

	mockRepo.On("GetByStatus", "sent", 1000).Return(messages, nil)

	service := NewMessageService(mockRepo, stopChan)
	result, err := service.Retrieve("sent")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "+905551111111", result[0].PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestStartProcess(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool)
	service := NewMessageService(mockRepo, stopChan)

	messages := []entity.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Test Process", Status: "pending"},
	}

	mockRepo.On("GetByStatus", "pending", mock.AnythingOfType("int")).
		Return(messages, nil).Maybe()

	mockRepo.On("Update", uint(1), "sent").Return(nil).Maybe()

	go service.StartProcess()
	time.Sleep(3 * time.Second)
	service.StopProcess()

	time.Sleep(500 * time.Millisecond)

	mockRepo.AssertExpectations(t)
}

func TestSendToWebhook(t *testing.T) {
	os.Setenv("AUTH_KEY", "test-auth-key")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-ins-auth-key") != "test-auth-key" || r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		var p model.Payload
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(model.Response{MessageID: "test-id", Message: "OK"})
	}))
	defer ts.Close()
	os.Setenv("WEBHOOK_URL", ts.URL)

	svc := NewMessageService(nil, nil)
	resp, err := svc.sendToWebhook(model.Payload{
		To:      "+123456789",
		Content: "Hello",
	})
	assert.NoError(t, err)
	assert.Equal(t, "test-id", resp.MessageID)
}

func TestProcessMessages(t *testing.T) {
	os.Setenv("AUTH_KEY", "test-auth-key")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(model.Response{MessageID: "test-id", Message: "OK"})
	}))
	defer ts.Close()
	os.Setenv("WEBHOOK_URL", ts.URL)

	mockRepo := new(MockMessageRepo)
	pendingStatus := entity.StatusPending
	msgs := []entity.Message{
		{ID: 1, PhoneNumber: "+123456789", Content: "Test message", Status: pendingStatus, CreatedAt: time.Now()},
	}
	mockRepo.On("GetByStatus", pendingStatus, messageCountPerMinute).Return(msgs, nil)
	mockRepo.On("Update", uint(1), entity.StatusSent).Return(nil)

	stopChan := make(chan bool)
	svc := NewMessageService(mockRepo, stopChan)
	err := svc.process()
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessMessages_GetByStatusError(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	pendingStatus := entity.StatusPending
	mockRepo.On("GetByStatus", pendingStatus, messageCountPerMinute).Return([]entity.Message{}, errors.New("db error"))
	stopChan := make(chan bool)
	svc := NewMessageService(mockRepo, stopChan)
	err := svc.process()
	assert.Error(t, err)
	assert.Equal(t, "error occurred when getting messages", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestProcessMessages_SendToWebhookError(t *testing.T) {
	os.Setenv("AUTH_KEY", "test-auth-key")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "internal error", http.StatusInternalServerError)
	}))
	defer ts.Close()
	os.Setenv("WEBHOOK_URL", ts.URL)

	mockRepo := new(MockMessageRepo)
	pendingStatus := entity.StatusPending
	msgs := []entity.Message{
		{ID: 2, PhoneNumber: "+123456789", Content: "Test message", Status: pendingStatus, CreatedAt: time.Now()},
	}
	mockRepo.On("GetByStatus", pendingStatus, messageCountPerMinute).Return(msgs, nil)
	// In this test, Update is not expected because sendToWebhook fails.
	stopChan := make(chan bool)
	svc := NewMessageService(mockRepo, stopChan)
	err := svc.process()
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
