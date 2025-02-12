package handler

import (
	"encoding/json"
	"github.com/busragumusel/insider-case/internal/service"
	"net/http"
)

type MessageHandler struct {
	service service.MessageSvc
}

func NewMessageHandler(service service.MessageSvc) *MessageHandler {
	return &MessageHandler{service: service}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, response interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, `{"code":"ERROR","message":"Failed to encode response"}`, http.StatusInternalServerError)
	}
}

// StartProcess starts message processing
// @Summary Start message processing
// @Description Starts the background process that handles messages.
// @Tags Message
// @Produce json
// @Success 200 {object} APIResult
// @Router /start [get]
func (r *MessageHandler) StartProcess(w http.ResponseWriter, r2 *http.Request) {
	go r.service.StartProcess()
	writeJSONResponse(w, http.StatusOK, nil)
}

// StopProcess stops message processing
// @Summary Stop message processing
// @Description Stops the message processing Goroutine.
// @Tags Message
// @Produce json
// @Success 200 {object} APIResult
// @Router /stop [get]
func (r *MessageHandler) StopProcess(w http.ResponseWriter, r2 *http.Request) {
	r.service.StopProcess()
	writeJSONResponse(w, http.StatusOK, nil)
}

// Retrieve fetches all sent messages
// @Summary Retrieve sent messages
// @Description Fetches all sent messages from the database.
// @Tags Message
// @Produce json
// @Success 200 {object} APIResult
// @Router /messages [get]
func (r *MessageHandler) Retrieve(w http.ResponseWriter, req *http.Request) {
	status := req.URL.Query().Get("status") // ✅ Get `status` from query params

	messages, err := r.service.Retrieve(status)
	if err != nil {
		writeJSONResponse(w, http.StatusInternalServerError, APIError{
			Message: "Failed to fetch messages",
		})
		return
	}

	writeJSONResponse(w, http.StatusOK, APIResult{
		Data: messages,
	})
}
