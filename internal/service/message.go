package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/busragumusel/insider-case/internal/entity"
	"github.com/busragumusel/insider-case/internal/model"
	"github.com/busragumusel/insider-case/internal/repository"
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
	StartProcess()
	StopProcess()
	Retrieve(status string) ([]entity.Message, error)
}

type MessageService struct {
	repo     repository.MessageRepo
	stopChan chan bool
}

func NewMessageService(
	repo repository.MessageRepo,
	stopChan chan bool,
) *MessageService {
	return &MessageService{
		repo,
		stopChan,
	}
}

func (s *MessageService) StartProcess() {
	ticker := time.NewTicker(processTimeRange * time.Minute)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				if err := s.process(); err != nil {
					log.Println("Error processing messages:", err)
				}
			case <-s.stopChan:
				log.Println("Stopping process...")
				return
			}
		}
	}()
}

func (s *MessageService) StopProcess() {
	close(s.stopChan)
	log.Println("Process stopped.")
}

func (s *MessageService) Retrieve(status string) ([]entity.Message, error) {
	messages, err := s.repo.GetByStatus(status, 1000)
	if err != nil {
		return nil, errors.New("failed to retrieve sent messages")
	}

	return messages, nil
}

func (s *MessageService) process() error {
	messages, err := s.repo.GetByStatus(entity.StatusPending, messageCountPerMinute)
	if err != nil {
		return errors.New("error occurred when getting messages")
	}

	if len(messages) == 0 {
		log.Println("No pending messages to send.")
		return nil
	}

	for _, msg := range messages {
		res, err := s.sendToWebhook(model.Payload{
			To:      msg.PhoneNumber,
			Content: msg.Content,
		})
		if err != nil {
			log.Println("Failed to send message:", err)
			continue
		}

		log.Printf("messageID: %s", res.MessageID)

		err = s.repo.Update(msg.ID, entity.StatusSent)
		if err != nil {
			log.Printf("Failed to update message with ID %d: %v", msg.ID, err)
		}
		log.Printf("Message sent! ID: %d\n", msg.ID)
	}

	return nil
}

func (s *MessageService) sendToWebhook(payload model.Payload) (model.Response, error) {
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

	return response, nil
}
