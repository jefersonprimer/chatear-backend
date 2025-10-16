package usecases

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// MockUserRepository is a mock implementation of repositories.UserRepository
type MockUserRepository struct {
	CreateFunc           func(ctx context.Context, user *entities.User) error
	GetByIDFunc          func(ctx context.Context, id uuid.UUID) (*entities.User, error)
	GetByEmailFunc       func(ctx context.Context, email string) (*entities.User, error)
	UpdateFunc           func(ctx context.Context, user *entities.User) error
	DeleteFunc           func(ctx context.Context, id uuid.UUID) error
	GetDeletedUsersFunc  func(ctx context.Context, limit, offset int) ([]*entities.User, error)
	GetByEmailVerifiedFunc func(ctx context.Context, verified bool, limit, offset int) ([]*entities.User, error)
	SearchByNameFunc     func(ctx context.Context, name string, limit, offset int) ([]*entities.User, error)
}

// Create implements repositories.UserRepository
func (m *MockUserRepository) Create(ctx context.Context, user *entities.User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}
	return errors.New("Create not implemented")
}

// GetByID implements repositories.UserRepository
func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}
	return nil, errors.New("GetByID not implemented")
}

// GetByEmail implements repositories.UserRepository
func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}
	return nil, errors.New("GetByEmail not implemented")
}

// Update implements repositories.UserRepository
func (m *MockUserRepository) Update(ctx context.Context, user *entities.User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}
	return errors.New("Update not implemented")
}

// Delete implements repositories.UserRepository
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}
	return errors.New("Delete not implemented")
}

// GetDeletedUsers implements repositories.UserRepository
func (m *MockUserRepository) GetDeletedUsers(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	if m.GetDeletedUsersFunc != nil {
		return m.GetDeletedUsersFunc(ctx, limit, offset)
	}
	return nil, errors.New("GetDeletedUsers not implemented")
}

// GetByEmailVerified implements repositories.UserRepository
func (m *MockUserRepository) GetByEmailVerified(ctx context.Context, verified bool, limit, offset int) ([]*entities.User, error) {
	if m.GetByEmailVerifiedFunc != nil {
		return m.GetByEmailVerifiedFunc(ctx, verified, limit, offset)
	}
	return nil, errors.New("GetByEmailVerified not implemented")
}

// SearchByName implements repositories.UserRepository
func (m *MockUserRepository) SearchByName(ctx context.Context, name string, limit, offset int) ([]*entities.User, error) {
	if m.SearchByNameFunc != nil {
		return m.SearchByNameFunc(ctx, name, limit, offset)
	}
	return nil, errors.New("SearchByName not implemented")
}

func TestGetUser(t *testing.T) {
	userID := uuid.New()
	expectedUser := &entities.User{
		ID:    userID,
		Email: "test@example.com",
		Name:  "Test User",
	}

	mockRepo := &MockUserRepository{
		GetByIDFunc: func(ctx context.Context, id uuid.UUID) (*entities.User, error) {
			if id == userID {
				return expectedUser, nil
			}
			return nil, errors.New("user not found")
		},
	}

	useCases := NewUserUseCases(mockRepo)

	// Test case 1: User found
	user, err := useCases.GetUser(context.Background(), userID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if user.(*entities.User).ID != expectedUser.ID {
		t.Errorf("Expected user ID %s, got %s", expectedUser.ID, user.(*entities.User).ID)
	}

	// Test case 2: User not found
	nonExistentUserID := uuid.New()
	_, err = useCases.GetUser(context.Background(), nonExistentUserID)
	if err == nil {
		t.Fatal("Expected an error, got nil")
	}
	if err.Error() != "user not found" {
		t.Errorf("Expected error 'user not found', got %v", err)
	}
}

func TestCreateUser(t *testing.T) {
	// Stub test for CreateUser
	t.Skip("Not implemented yet")
}

func TestUpdateUser(t *testing.T) {
	// Stub test for UpdateUser
	t.Skip("Not implemented yet")
}

func TestDeleteUser(t *testing.T) {
	// Stub test for DeleteUser
	t.Skip("Not implemented yet")
}