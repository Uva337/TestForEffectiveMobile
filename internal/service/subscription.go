package service

import (
	"context"
	"fmt"
	"time"

	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository/postgres"
	"github.com/google/uuid"
)


type SubscriptionRepository interface {
	Create(ctx context.Context, sub *models.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error)
	Update(ctx context.Context, sub *models.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetSummary(ctx context.Context, filter postgres.GetSummaryFilter) (int, error)
}


type SubscriptionService struct {
	repo SubscriptionRepository
}


func NewSubscriptionService(repo SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		repo: repo,
	}
}


type CreateSubscriptionDTO struct {
	UserID      uuid.UUID  `json:"user_id" validate:"required"`                 
	ServiceName string     `json:"service_name" validate:"required,min=2,max=100"` 
	Price       int        `json:"price" validate:"required,gt=0"`              
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}


func (s *SubscriptionService) Create(ctx context.Context, dto CreateSubscriptionDTO) (*models.Subscription, error) {
	sub := &models.Subscription{
		ID:          uuid.New(), 
		UserID:      dto.UserID,
		ServiceName: dto.ServiceName,
		Price:       dto.Price,
		StartDate:   dto.StartDate,
		EndDate:     dto.EndDate,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, fmt.Errorf("не удалось создать подписку: %w", err)
	}

	return sub, nil
}


func (s *SubscriptionService) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}


type UpdateSubscriptionDTO struct {
	ServiceName string     `json:"service_name" validate:"required,min=2,max=100"`
	Price       int        `json:"price" validate:"required,gt=0"`
	StartDate   time.Time  `json:"start_date" validate:"required"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}


func (s *SubscriptionService) Update(ctx context.Context, id uuid.UUID, dto UpdateSubscriptionDTO) error {
	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Обновляем поля
	sub.ServiceName = dto.ServiceName
	sub.Price = dto.Price
	sub.StartDate = dto.StartDate
	sub.EndDate = dto.EndDate

	return s.repo.Update(ctx, sub)
}


func (s *SubscriptionService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}


func (s *SubscriptionService) GetSummary(ctx context.Context, filter postgres.GetSummaryFilter) (int, error) {
	return s.repo.GetSummary(ctx, filter)
}
