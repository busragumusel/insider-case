package service

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/busragumusel/insider-case/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockMessageRepo struct {
	mock.Mock
}

func (m *MockMessageRepo) GetByStatus(ctx context.Context, status string, limit int) ([]entity.Message, error) {
	args := m.Called(ctx, status, limit)
	return args.Get(0).([]entity.Message), args.Error(1)
}

func (m *MockMessageRepo) Update(ctx context.Context, id uint, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func setupRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func TestRetrieveMessages(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool)
	ctx := context.Background()
	redisClient := setupRedisClient()

	messages := []entity.Message{
		{ID: 1, PhoneNumber: "+905551111111", Content: "Test Message", Status: "sent"},
	}

	mockRepo.On("GetByStatus", ctx, "sent", 1000).Return(messages, nil)

	service := NewMessageService(mockRepo, stopChan, redisClient)
	result, err := service.Retrieve(ctx, "sent")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "+905551111111", result[0].PhoneNumber)
	mockRepo.AssertExpectations(t)
}

func TestProcessMessages(t *testing.T) {
	os.Setenv("WEBHOOK_URL", "http://example.com/webhook")
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
		json.NewEncoder(w).Encode(model.Response{MessageID: "test-id", Message: "OK"})
	}))
	defer ts.Close()
	os.Setenv("WEBHOOK_URL", ts.URL)

	mockRepo := new(MockMessageRepo)
	ctx := context.Background()
	redisClient := setupRedisClient()
	pendingStatus := entity.StatusPending
	msgs := []entity.Message{
		{ID: 1, PhoneNumber: "+123456789", Content: "Test message", Status: pendingStatus, CreatedAt: time.Now()},
	}
	mockRepo.On("GetByStatus", ctx, pendingStatus, messageCountPerMinute).Return(msgs, nil)
	mockRepo.On("Update", ctx, uint(1), entity.StatusSent).Return(nil)

	stopChan := make(chan bool)
	svc := NewMessageService(mockRepo, stopChan, redisClient)
	err := svc.process(ctx)
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessMessages_GetByStatusError(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	ctx := context.Background()
	redisClient := setupRedisClient()
	pendingStatus := entity.StatusPending
	mockRepo.On("GetByStatus", ctx, pendingStatus, messageCountPerMinute).Return([]entity.Message{}, errors.New("db error"))
	stopChan := make(chan bool)
	svc := NewMessageService(mockRepo, stopChan, redisClient)
	err := svc.process(ctx)
	assert.Error(t, err)
	assert.Equal(t, "error occurred when getting messages", err.Error())
	mockRepo.AssertExpectations(t)
}
