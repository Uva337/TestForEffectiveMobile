package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository/postgres"
	"effective-mobile-task/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10" 
	"github.com/google/uuid"
)


type SubscriptionService interface {
	Create(ctx context.Context, dto service.CreateSubscriptionDTO) (*models.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, dto service.UpdateSubscriptionDTO) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSummary(ctx context.Context, filter postgres.GetSummaryFilter) (int, error)
}


type Handler struct {
	service  SubscriptionService
	log      *slog.Logger
	validate *validator.Validate 
}


func NewHandler(service SubscriptionService, log *slog.Logger) *Handler {
	return &Handler{
		service:  service,
		log:      log,
		validate: validator.New(), 
	}
}

// CreateSubscription обрабатывает запрос на создание новой подписки.
// @Summary Create a new subscription
// @Description Add a new subscription to the database
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   subscription  body      service.CreateSubscriptionDTO  true  "Subscription Info"
// @Success 201           {object}  models.Subscription
// @Failure 400           {string}  string "Неверный формат JSON или неверные данные"
// @Failure 500           {string}  string "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var dto service.CreateSubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		h.log.Warn("не удалось декодировать тело запроса", "error", err)
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}


	if err := h.validate.Struct(dto); err != nil {
		h.log.Warn("неверные данные", "error", err)

		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	sub, err := h.service.Create(r.Context(), dto)
	if err != nil {
		h.log.Error("не удалось создать подписку", "error", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, sub)
}

// GetSubscriptionByID обрабатывает запрос на получение подписки по ID.
// @Summary Get a subscription by ID
// @Description Get details of a specific subscription
// @Tags subscriptions
// @Produce  json
// @Param   id   path      string  true  "Subscription ID"
// @Success 200  {object}  models.Subscription
// @Failure 400  {string}  string "Неверный формат ID"
// @Failure 404  {string}  string "Подписка не найдена"
// @Failure 500  {string}  string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	sub, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		h.log.Error("не удалось получить подписку", "id", id, "error", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, sub)
}

// UpdateSubscription обрабатывает запрос на обновление подписки.
// @Summary Update an existing subscription
// @Description Update details of an existing subscription by its ID
// @Tags subscriptions
// @Accept  json
// @Produce  json
// @Param   id            path      string                       true  "Subscription ID"
// @Param   subscription  body      service.UpdateSubscriptionDTO  true  "Subscription data to update"
// @Success 200           {string}  string "OK"
// @Failure 400           {string}  string "Неверный формат JSON или неверные данные"
// @Failure 404           {string}  string "Подписка не найдена"
// @Failure 500           {string}  string "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	var dto service.UpdateSubscriptionDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}


	if err := h.validate.Struct(dto); err != nil {
		h.log.Warn("неверные данные", "error", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}


	err = h.service.Update(r.Context(), id, dto)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		h.log.Error("не удалось обновить подписку", "id", id, "error", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// DeleteSubscription обрабатывает запрос на удаление подписки.
// @Summary Delete a subscription
// @Description Delete a subscription by its ID
// @Tags subscriptions
// @Produce  json
// @Param   id   path      string  true  "Subscription ID"
// @Success 204  {string}  string "No Content"
// @Failure 400  {string}  string "Invalid ID format"
// @Failure 404  {string}  string "Subscription not found"
// @Failure 500  {string}  string "Internal server error"
// @Router /subscriptions/{id} [delete]
func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Неверный формат ID", http.StatusBadRequest)
		return
	}

	err = h.service.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, postgres.ErrNotFound) {
			http.Error(w, "Подписка не найдена", http.StatusNotFound)
			return
		}
		h.log.Error("не удалось удалить подписку", "id", id, "error", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetSummary обрабатывает запрос на получение суммарной стоимости.
// @Summary Get summary price of subscriptions
// @Description Calculates the total price of subscriptions based on optional filters
// @Tags subscriptions
// @Produce  json
// @Param   user_id       query     string  false  "Filter by User ID (UUID format)"
// @Param   service_name  query     string  false  "Filter by Service Name"
// @Param   start_date    query     string  false  "Filter by start date (YYYY-MM-DD)"
// @Param   end_date      query     string  false  "Filter by end date (YYYY-MM-DD)"
// @Success 200           {object}  map[string]int
// @Failure 400           {string}  string "Invalid filter format"
// @Failure 500           {string}  string "Internal server error"
// @Router /subscriptions/summary [get]
func (h *Handler) GetSummary(w http.ResponseWriter, r *http.Request) {
	var filter postgres.GetSummaryFilter
	q := r.URL.Query()
	const layout = "2006-01-02" // Формат для парсинга YYYY-MM-DD

	if userIDStr := q.Get("user_id"); userIDStr != "" {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			http.Error(w, "Неверный формат user_id", http.StatusBadRequest)
			return
		}
		filter.UserID = &userID
	}

	if serviceName := q.Get("service_name"); serviceName != "" {
		filter.ServiceName = &serviceName
	}

	if startDateStr := q.Get("start_date"); startDateStr != "" {
		startDate, err := time.Parse(layout, startDateStr)
		if err != nil {
			http.Error(w, "Неверный формат start_date, используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		filter.StartDate = &startDate
	}

	if endDateStr := q.Get("end_date"); endDateStr != "" {
		endDate, err := time.Parse(layout, endDateStr)
		if err != nil {
			http.Error(w, "Неверный формат end_date, используйте YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		filter.EndDate = &endDate
	}

	total, err := h.service.GetSummary(r.Context(), filter)
	if err != nil {
		h.log.Error("не удалось получить сводку", "error", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"total_price": total})
}


func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}


