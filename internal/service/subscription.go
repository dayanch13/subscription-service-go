package service

import (
	"github.com/google/uuid"
	"log"
	"subscription-service-go/internal/model"
	"subscription-service-go/internal/repository/postgres"
)

type SubscriptionService struct {
	repo *postgres.SubscriptionRepository
}

func NewSubscriptionService(repo *postgres.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{repo: repo}
}

func (s *SubscriptionService) CreateSubscription(createReq *model.SubscriptionCreate) (*model.Subscription, error) {
	log.Printf("Creating subscription for user: %s", createReq.UserID)

	subscription := &model.Subscription{
		ServiceName: createReq.ServiceName,
		Price:       createReq.Price,
		UserID:      createReq.UserID,
		StartDate:   createReq.StartDate,
		EndDate:     createReq.EndDate,
	}

	if err := s.repo.Create(subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

func (s *SubscriptionService) GetSubscription(id int) (*model.Subscription, error) {
	log.Printf("Getting subscription with ID: %d", id)
	return s.repo.GetByID(id)
}

func (s *SubscriptionService) GetUserSubscriptions(userID string) ([]model.Subscription, error) {
	log.Printf("Getting subscriptions for user: %s", userID)

	uuid, err := parseUUID(userID)
	if err != nil {
		return nil, err
	}

	return s.repo.GetByUserID(uuid)
}

func (s *SubscriptionService) GetAllSubscriptions() ([]model.Subscription, error) {
	log.Println("Getting all subscriptions")
	return s.repo.GetAll()
}

func (s *SubscriptionService) UpdateSubscription(id int, updateReq *model.SubscriptionUpdate) error {
	log.Printf("Updating subscription with ID: %d", id)
	return s.repo.Update(id, updateReq)
}

func (s *SubscriptionService) DeleteSubscription(id int) error {
	log.Printf("Deleting subscription with ID: %d", id)
	return s.repo.Delete(id)
}

func (s *SubscriptionService) CalculateCost(req *model.CostRequest) (int, error) {
	log.Printf("Calculating cost for period: %s to %s", req.StartPeriod, req.EndPeriod)
	return s.repo.CalculateCost(req)
}

func parseUUID(uuidStr string) (uuid.UUID, error) {
	return uuid.Parse(uuidStr)
}
