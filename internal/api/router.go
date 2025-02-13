package api

import (
	"github.com/busragumusel/insider-case/internal/handler"
	"github.com/busragumusel/insider-case/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type API struct {
	db             *gorm.DB
	redisClient    *redis.Client
	messageService *service.MessageService
}

func NewAPI(
	db *gorm.DB,
	redisClient *redis.Client,
	messageService *service.MessageService,
) *API {
	return &API{
		db:             db,
		redisClient:    redisClient,
		messageService: messageService,
	}
}

func (r *API) RegisterRoutes(router *chi.Mux) {
	messageHandler := handler.NewMessageHandler(r.messageService)

	router.Get("/start", messageHandler.StartProcess)
	router.Get("/stop", messageHandler.StopProcess)
	router.Get("/messages", messageHandler.Retrieve)
}
