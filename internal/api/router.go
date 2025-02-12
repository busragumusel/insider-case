package api

import (
	"github.com/busragumusel/insider-case/internal/handler"
	"github.com/busragumusel/insider-case/internal/repository"
	"github.com/busragumusel/insider-case/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type API struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewAPI(
	db *gorm.DB,
	redisClient *redis.Client,
) *API {
	return &API{
		db:          db,
		redisClient: redisClient,
	}
}

func (r *API) RegisterRoutes(router *chi.Mux) {
	messageRepository := repository.NewMessageRepository(r.db)

	stopChan := make(chan bool)
	messageService := service.NewMessageService(messageRepository, stopChan, r.redisClient)

	messageHandler := handler.NewMessageHandler(messageService)

	router.Get("/start", messageHandler.StartProcess)
	router.Get("/stop", messageHandler.StopProcess)
	router.Get("/messages", messageHandler.Retrieve)
}
