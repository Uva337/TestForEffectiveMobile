package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"effective-mobile-task/internal/models"
	"effective-mobile-task/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock" 
)


type MockRepository struct {
	mock.Mock 
}



func (m *MockRepository) Create(ctx context.Context, sub *models.Subscription) error {
	
	args := m.Called(ctx, sub)
	
	return args.Error(0)
}

func (m *MockRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	args := m.Called(ctx, id)
	
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Subscription), args.Error(1)
}

func (m *MockRepository) Update(ctx context.Context, sub *models.Subscription) error {
	args := m.Called(ctx, sub)
	return args.Error(0)
}

func (m *MockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRepository) GetSummary(ctx context.Context, filter postgres.GetSummaryFilter) (int, error) {
	args := m.Called(ctx, filter)
	
	return args.Int(0), args.Error(1)
}



func TestSubscriptionService_Create_Success(t *testing.T) {
	
	mockRepo := new(MockRepository)
	service := NewSubscriptionService(mockRepo)

	dto := CreateSubscriptionDTO{
		UserID:      uuid.New(),
		ServiceName: "Test Service",
		Price:       100,
		StartDate:   time.Now(),
	}

	
	mockRepo.On("Create", mock.Anything, mock.AnythingOfType("*models.Subscription")).Return(nil)

	
	sub, err := service.Create(context.Background(), dto)


	assert.NoError(t, err)                   
	assert.NotNil(t, sub)                
	assert.Equal(t, dto.ServiceName, sub.ServiceName) 
	assert.NotEqual(t, uuid.Nil, sub.ID)    
	mockRepo.AssertExpectations(t)           
}


func TestSubscriptionService_GetByID_NotFound(t *testing.T) {

	mockRepo := new(MockRepository)
	service := NewSubscriptionService(mockRepo)
	
	testID := uuid.New()
	expectedError := errors.New("not found") 

	
	mockRepo.On("GetByID", mock.Anything, testID).Return(nil, expectedError)


	sub, err := service.GetByID(context.Background(), testID)

	
	assert.Error(t, err)                
	assert.Nil(t, sub)                  
	assert.Equal(t, expectedError, err) 
	mockRepo.AssertExpectations(t)
}


