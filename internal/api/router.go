package api

import (
	"github.com/busragumusel/insider-case/internal/handler"
	"github.com/busragumusel/insider-case/internal/repository"
	"github.com/busragumusel/insider-case/internal/service"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type API struct {
	db *gorm.DB
}

func NewAPI(
	db *gorm.DB,
) *API {
	return &API{db: db}
}

func (r *API) RegisterRoutes(router *chi.Mux) {
	messageRepository := repository.NewMessageRepository(r.db)

	stopChan := make(chan bool)
	messageService := service.NewMessageService(messageRepository, stopChan)

	messageHandler := handler.NewMessageHandler(messageService)

	router.Get("/start", messageHandler.StartProcess)
	router.Get("/stop", messageHandler.StopProcess)
	router.Get("/messages", messageHandler.Retrieve)
}
