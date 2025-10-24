package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"effective-mobile-task/internal/models"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("subscription not found")

type SubscriptionRepository struct {
	db *pgxpool.Pool
	
	sqb sq.StatementBuilderType
}


func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{
		db:  db,
		sqb: sq.StatementBuilder.PlaceholderFormat(sq.Dollar), 
	}
}


func (r *SubscriptionRepository) Create(ctx context.Context, sub *models.Subscription) error {
	sql, args, err := r.sqb.Insert("subscriptions").
		Columns("id", "user_id", "service_name", "price", "start_date", "end_date").
		Values(sub.ID, sub.UserID, sub.ServiceName, sub.Price, sub.StartDate, sub.EndDate).
		ToSql()
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Create - ToSql: %w", err)
	}

	_, err = r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Create - Exec: %w", err)
	}

	return nil
}


func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Subscription, error) {
	sql, args, err := r.sqb.Select("id", "user_id", "service_name", "price", "start_date", "end_date").
		From("subscriptions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("SubscriptionRepository.GetByID - ToSql: %w", err)
	}

	var sub models.Subscription
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&sub.ID,
		&sub.UserID,
		&sub.ServiceName,
		&sub.Price,
		&sub.StartDate,
		&sub.EndDate,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("SubscriptionRepository.GetByID - Scan: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, sub *models.Subscription) error {
	sql, args, err := r.sqb.Update("subscriptions").
		Set("service_name", sub.ServiceName).
		Set("price", sub.Price).
		Set("start_date", sub.StartDate).
		Set("end_date", sub.EndDate).
		Where(sq.Eq{"id": sub.ID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Update - ToSql: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Update - Exec: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	sql, args, err := r.sqb.Delete("subscriptions").
		Where(sq.Eq{"id": id}).
		ToSql()
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Delete - ToSql: %w", err)
	}

	res, err := r.db.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("SubscriptionRepository.Delete - Exec: %w", err)
	}

	if res.RowsAffected() == 0 {
		return ErrNotFound
	}

	return nil
}

type GetSummaryFilter struct {
	UserID      *uuid.UUID
	ServiceName *string
	StartDate   *time.Time 
	EndDate     *time.Time 
}

func (r *SubscriptionRepository) GetSummary(ctx context.Context, filter GetSummaryFilter) (int, error) {
	queryBuilder := r.sqb.Select("COALESCE(SUM(price), 0)").From("subscriptions")

	if filter.UserID != nil {
		queryBuilder = queryBuilder.Where(sq.Eq{"user_id": *filter.UserID})
	}
	if filter.ServiceName != nil {
		queryBuilder = queryBuilder.Where(sq.Eq{"service_name": *filter.ServiceName})
	}
	if filter.StartDate != nil {
		queryBuilder = queryBuilder.Where(sq.GtOrEq{"start_date": *filter.StartDate})
	}
	if filter.EndDate != nil {
		queryBuilder = queryBuilder.Where(sq.LtOrEq{"start_date": *filter.EndDate})
	}

	sql, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("SubscriptionRepository.GetSummary - ToSql: %w", err)
	}

	var total int
	err = r.db.QueryRow(ctx, sql, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("SubscriptionRepository.GetSummary - Scan: %w", err)
	}

	return total, nil
}
