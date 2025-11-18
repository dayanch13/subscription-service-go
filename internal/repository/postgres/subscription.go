package postgres

import (
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"log"
	"subscription-service-go/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	query := `
        INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_at, updated_at
    `

	err := r.db.QueryRow(
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(&sub.ID, &sub.CreatedAt, &sub.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	log.Printf("Created subscription with ID: %d", sub.ID)
	return nil
}

func (r *SubscriptionRepository) GetByID(id int) (*model.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE id = $1`

	var sub model.Subscription
	err := r.db.QueryRow(query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetByUserID(userID uuid.UUID) ([]model.Subscription, error) {
	query := `SELECT * FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	log.Printf("Found %d subscriptions for user %s", len(subscriptions), userID)
	return subscriptions, nil
}

func (r *SubscriptionRepository) GetAll() ([]model.Subscription, error) {
	query := `SELECT * FROM subscriptions ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []model.Subscription
	for rows.Next() {
		var sub model.Subscription
		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, sub)
	}

	log.Printf("Found %d total subscriptions", len(subscriptions))
	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(id int, update *model.SubscriptionUpdate) error {
	query := `
        UPDATE subscriptions 
        SET service_name = COALESCE($1, service_name),
            price = COALESCE($2, price),
            start_date = COALESCE($3, start_date),
            end_date = COALESCE($4, end_date),
            updated_at = CURRENT_TIMESTAMP
        WHERE id = $5
    `

	result, err := r.db.Exec(
		query,
		update.ServiceName,
		update.Price,
		update.StartDate,
		update.EndDate,
		id,
	)

	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Updated subscription with ID: %d", id)
	return nil
}

func (r *SubscriptionRepository) Delete(id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	log.Printf("Deleted subscription with ID: %d", id)
	return nil
}

func (r *SubscriptionRepository) CalculateCost(req *model.CostRequest) (int, error) {
	query := `
        SELECT COALESCE(SUM(price), 0) 
        FROM subscriptions 
        WHERE (end_date IS NULL OR end_date >= $1)
        AND start_date <= $2
    `

	args := []interface{}{req.StartPeriod, req.EndPeriod}
	argIndex := 3

	if req.UserID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argIndex)
		args = append(args, *req.UserID)
		argIndex++
	}

	if req.ServiceName != nil {
		query += fmt.Sprintf(" AND service_name ILIKE $%d", argIndex)
		args = append(args, "%"+*req.ServiceName+"%")
	}

	var totalCost int
	err := r.db.QueryRow(query, args...).Scan(&totalCost)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate cost: %w", err)
	}

	log.Printf("Calculated total cost: %d for period %s-%s", totalCost, req.StartPeriod, req.EndPeriod)
	return totalCost, nil
}
