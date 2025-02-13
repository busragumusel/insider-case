package service

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/busragumusel/insider-case/internal/entity"
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

func setupRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})
}

func TestStartProcess(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool, 1)
	ctx := context.Background()
	redisClient := setupRedisClient()
	var mu sync.Mutex

	service := NewMessageService(mockRepo, stopChan, redisClient, &mu, false)

	// Start the process
	service.StartProcess(ctx)

	// Wait for a short time to simulate processing
	time.Sleep(100 * time.Millisecond)

	assert.True(t, service.running, "Service should be running after StartProcess")

	// Stop process to clean up goroutine
	service.StopProcess()
}

func TestStopProcess(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool, 1)
	ctx := context.Background()
	redisClient := setupRedisClient()
	var mu sync.Mutex

	service := NewMessageService(mockRepo, stopChan, redisClient, &mu, false)

	// Start the process
	service.StartProcess(ctx)

	// Ensure the service is running
	time.Sleep(100 * time.Millisecond)
	assert.True(t, service.running, "Service should be running before StopProcess")

	// Stop the process
	service.StopProcess()

	// Allow time for stop to take effect
	time.Sleep(100 * time.Millisecond)
	assert.False(t, service.running, "Service should not be running after StopProcess")
}

func TestMultipleStartProcessCalls(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool, 1)
	ctx := context.Background()
	redisClient := setupRedisClient()
	var mu sync.Mutex

	service := NewMessageService(mockRepo, stopChan, redisClient, &mu, false)

	// Start the process
	service.StartProcess(ctx)
	time.Sleep(50 * time.Millisecond)

	// Start process again - it should not create a duplicate
	service.StartProcess(ctx)
	time.Sleep(50 * time.Millisecond)

	assert.True(t, service.running, "Service should still be running after multiple StartProcess calls")

	// Stop process to clean up goroutine
	service.StopProcess()
	time.Sleep(100 * time.Millisecond)
	assert.False(t, service.running, "Service should stop after StopProcess")
}

func TestProcessStopsOnStopChan(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	stopChan := make(chan bool, 1)
	ctx := context.Background()
	redisClient := setupRedisClient()
	var mu sync.Mutex

	service := NewMessageService(mockRepo, stopChan, redisClient, &mu, false)

	// Start the process
	service.StartProcess(ctx)
	time.Sleep(100 * time.Millisecond)

	// Simulate stop signal
	stopChan <- true

	// Allow time for goroutine to stop
	time.Sleep(100 * time.Millisecond)

	assert.False(t, service.running, "Service should not be running after receiving stop signal")
}

func TestProcessHandlesDBErrorGracefully(t *testing.T) {
	mockRepo := new(MockMessageRepo)
	ctx := context.Background()
	redisClient := setupRedisClient()
	var mu sync.Mutex

	mockRepo.On("GetByStatus", ctx, entity.StatusPending, messageCountPerMinute).Return([]entity.Message{}, errors.New("DB error"))

	stopChan := make(chan bool, 1)
	service := NewMessageService(mockRepo, stopChan, redisClient, &mu, false)

	// Test process handling DB error
	err := service.process(ctx)
	assert.Error(t, err)
	assert.Equal(t, "error occurred when getting messages", err.Error())
	mockRepo.AssertExpectations(t)
}
