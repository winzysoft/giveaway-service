package giveaway

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"giveaway-service/internal/domain"
	"giveaway-service/internal/dto"
	"giveaway-service/internal/service/giveaway"
	"net/http"
)

type Handler struct {
	ctx             context.Context
	giveawayService giveaway.Service
}

func New(
	giveawayService *giveaway.Service,
) *Handler {
	return &Handler{
		giveawayService: *giveawayService,
	}
}

func (h *Handler) HandleFindGiveaway(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(respWriter, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 1. Извлечение ID из URL
	idStr := req.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(respWriter, "Invalid giveaway ID", http.StatusBadRequest)
		return
	}

	// 2. Вызов сервиса
	resp, err := h.giveawayService.Find(req.Context(), id)
	if err != nil {
		var httpErr *domain.HTTPError
		if errors.As(err, &httpErr) {
			http.Error(respWriter, httpErr.Message, httpErr.Code)
		} else {
			http.Error(respWriter, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	// 3. Отправка ответа
	respWriter.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(respWriter).Encode(resp)
	if err != nil {
		http.Error(respWriter, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) HandleSaveGiveaway(respWriter http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(respWriter, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Проверяем Content-Type
	if req.Header.Get("Content-Type") != "application/json" {
		http.Error(respWriter, "expected content-type application/json", http.StatusBadRequest)
		return
	}

	// Ограничиваем размер тела запроса
	req.Body = http.MaxBytesReader(respWriter, req.Body, 1048576) // 1MB

	// Декодируем JSON
	var request dto.GiveawaySaveRequest
	dec := json.NewDecoder(req.Body)
	dec.DisallowUnknownFields() // Запрещаем неизвестные поля

	if err := dec.Decode(&request); err != nil {
		http.Error(respWriter, err.Error(), http.StatusBadRequest)
		return
	}

	// Service layer call
	save, err := h.giveawayService.Save(req.Context(), request)
	if err != nil {
		var httpErr *domain.HTTPError
		if errors.As(err, &httpErr) {
			http.Error(respWriter, httpErr.Message, httpErr.Code)
		} else {
			http.Error(respWriter, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	respWriter.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(respWriter).Encode(save)
	if err != nil {
		http.Error(respWriter, "internal server error", http.StatusInternalServerError)
		return
	}
}
