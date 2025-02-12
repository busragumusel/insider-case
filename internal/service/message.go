package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/busragumusel/insider-case/internal/model"
	"github.com/busragumusel/insider-case/internal/repository"
	"github.com/go-redis/redis/v8"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	processTimeRange      = 2
	messageCountPerMinute = 2
)

type MessageSvc interface {
	StartProcess(ctx context.Context)
	StopProcess(ctx context.Context)
	Retrieve(ctx context.Context, status string) ([]entity.Message, error)
}

type MessageService struct {
	repo        repository.MessageRepo
	stopChan    chan bool
	redisClient *redis.Client
}

func NewMessageService(
	repo repository.MessageRepo,
	stopChan chan bool,
	redisClient *redis.Client,
) *MessageService {
	return &MessageService{
		repo,
		stopChan,
		redisClient,
	}
}

func (s *MessageService) StartProcess(ctx context.Context) {
	ticker := time.NewTicker(processTimeRange * time.Minute)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.process(ctx); err != nil {
					log.Println("Error processing messages:", err)
				}
			case <-s.stopChan:
				log.Println("Stopping process...")
				return
			}
		}
	}()
}

func (s *MessageService) StopProcess(_ context.Context) {
	close(s.stopChan)
	log.Println("Process stopped.")
}

func (s *MessageService) Retrieve(ctx context.Context, status string) ([]entity.Message, error) {
	messages, err := s.repo.GetByStatus(ctx, status, 1000)
	if err != nil {
		return nil, errors.New("failed to retrieve sent messages")
	}

	return messages, nil
}

func (s *MessageService) process(ctx context.Context) error {
	messages, err := s.repo.GetByStatus(ctx, entity.StatusPending, messageCountPerMinute)
	if err != nil {
		return errors.New("error occurred when getting messages")
	}

	if len(messages) == 0 {
		log.Println("No pending messages to send.")
		return nil
	}

	for _, msg := range messages {
		res, err := s.sendToWebhook(ctx, model.Payload{
			To:      msg.PhoneNumber,
			Content: msg.Content,
		})
		if err != nil {
			log.Println("Failed to send message:", err)
			continue
		}

		log.Printf("messageID: %s", res.MessageID)

		err = s.repo.Update(ctx, msg.ID, entity.StatusSent)
		if err != nil {
			log.Printf("Failed to update message with ID %d: %v", msg.ID, err)
		}
		log.Printf("Message sent! ID: %d\n", msg.ID)
	}

	return nil
}

func (s *MessageService) sendToWebhook(ctx context.Context, payload model.Payload) (model.Response, error) {
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return model.Response{}, err
	}

	req, err := http.NewRequest("POST", os.Getenv("WEBHOOK_URL"), bytes.NewBuffer(jsonData))
	if err != nil {
		return model.Response{}, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-ins-auth-key", os.Getenv("AUTH_KEY"))

	client := &http.Client{Timeout: 10 * time.Second}

	resp, err := client.Do(req)
	if err != nil {
		return model.Response{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		return model.Response{}, errors.New("failed to send request: " + resp.Status)
	}

	var response model.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return model.Response{}, err
	}

	s.saveToCache(ctx, response)

	return response, nil
}

func (s *MessageService) saveToCache(ctx context.Context, response model.Response) {
	sendingTime := time.Now().Format(time.RFC3339)
	key := "message_id:" + response.MessageID
	err := s.redisClient.HSet(ctx, key, map[string]interface{}{
		"sending_time": sendingTime,
	}).Err()
	if err != nil {
		log.Printf("failed to cache: %v", err)
	}
}
