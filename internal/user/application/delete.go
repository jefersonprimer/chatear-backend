package application

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

const maxDeletionsPerDay = 10

// DeleteUser is a use case for deleting a user.
type DeleteUser struct {
	UserRepository             domain.UserRepository
	UserDeletionRepository     domain.UserDeletionRepository
	DeletionCapacityRepository domain.DeletionCapacityRepository
	EventBus                   domain.EventBus
}

// NewDeleteUser creates a new DeleteUser use case.
func NewDeleteUser(userRepository domain.UserRepository, userDeletionRepository domain.UserDeletionRepository, deletionCapacityRepository domain.DeletionCapacityRepository, eventBus domain.EventBus) *DeleteUser {
	return &DeleteUser{
		UserRepository:             userRepository,
		UserDeletionRepository:     userDeletionRepository,
		DeletionCapacityRepository: deletionCapacityRepository,
		EventBus:                   eventBus,
	}
}

// Execute schedules a user for deletion.
func (uc *DeleteUser) Execute(ctx context.Context, id uuid.UUID) error {
	today := time.Now().Truncate(24 * time.Hour)
	capacity, err := uc.DeletionCapacityRepository.GetDeletionCapacity(ctx, today)
	if err != nil {
		// If there is no entry for today, we can assume the count is 0
		capacity = &domain.DeletionCapacity{Count: 0, MaxLimit: maxDeletionsPerDay}
	}

	if capacity.Count >= capacity.MaxLimit {
		return errors.New("deletion limit exceeded")
	}

	scheduledDate := time.Now().Add(90 * 24 * time.Hour)

	userDeletion := &domain.UserDeletion{
		ID:            uuid.New(),
		UserID:        id,
		ScheduledDate: scheduledDate,
		Status:        "scheduled",
	}

	if err := uc.UserDeletionRepository.CreateUserDeletion(ctx, userDeletion); err != nil {
		return err
	}

	userDeletionBytes, err := json.Marshal(userDeletion)
	if err != nil {
		return err
	}

	if err := uc.EventBus.Publish(ctx, &domain.Event{Subject: "user.delete", Data: userDeletionBytes}); err != nil {
		return err
	}

	return uc.DeletionCapacityRepository.IncrementDeletionCapacity(ctx, today)
}
