package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jefersonprimer/chatear-backend/domain/entities"
)

// MockPgxRow is a mock implementation of pgx.Row
type MockPgxRow struct {
	ScanFunc func(dest ...interface{}) error
}

// Scan implements pgx.Row
func (m *MockPgxRow) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

// MockPgxRows is a mock implementation of pgx.Rows
type MockPgxRows struct {
	NextFunc  func() bool
	ScanFunc  func(dest ...interface{}) error
	CloseFunc func()
	ErrFunc   func() error
}

// Next implements pgx.Rows
func (m *MockPgxRows) Next() bool {
	return m.NextFunc()
}

// Scan implements pgx.Rows
func (m *MockPgxRows) Scan(dest ...interface{}) error {
	return m.ScanFunc(dest...)
}

// Close implements pgx.Rows
func (m *MockPgxRows) Close() {
	m.CloseFunc()
}

// Err implements pgx.Rows
func (m *MockPgxRows) Err() error {
	return m.ErrFunc()
}

// CommandTag implements pgx.Rows
func (m *MockPgxRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag("MOCK")
}

// MockPgxCommandTag is a mock implementation of pgconn.CommandTag
type MockPgxCommandTag struct {
	RowsAffectedValue int64
}

// RowsAffected implements pgconn.CommandTag
func (m *MockPgxCommandTag) RowsAffected() int64 {
	return m.RowsAffectedValue
}

// String implements pgconn.CommandTag
func (m *MockPgxCommandTag) String() string {
	return ""
}

// MockPgxPool is a mock implementation of PgxPoolIface
type MockPgxPool struct {
	QueryRowFunc func(ctx context.Context, sql string, args ...interface{}) pgx.Row
	ExecFunc     func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	QueryFunc    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
}

// QueryRow implements PgxPoolIface
func (m *MockPgxPool) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(ctx, sql, args...)
	}
	return &MockPgxRow{ScanFunc: func(dest ...interface{}) error { return errors.New("QueryRow not implemented") }}
}

// Exec implements PgxPoolIface
func (m *MockPgxPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(ctx, sql, arguments...)
	}
	return pgconn.CommandTag{}, errors.New("Exec not implemented")
}

// Query implements PgxPoolIface
func (m *MockPgxPool) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(ctx, sql, args...)
	}
	return nil, errors.New("Query not implemented")
}

func TestCreateUser(t *testing.T) {
	userID := uuid.New()
	now := time.Now()
	user := &entities.User{
		ID:              userID,
		Name:            "Test User",
		Email:           "test@example.com",
		PasswordHash:    "hashedpassword",
		IsEmailVerified: true,
		AvatarURL:       &avatarURL,
	}

	mockPool := &MockPgxPool{
		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {
			return &MockPgxRow{
				ScanFunc: func(dest ...interface{}) error {
					*dest[0].(*time.Time) = now // created_at
					*dest[1].(*time.Time) = now // updated_at
					return nil
				},
			}
		},
	}

	repo := NewUserRepository(mockPool)

	err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("Create() failed: %v", err)
	}

	if user.CreatedAt.IsZero() || user.UpdatedAt.IsZero() {
		t.Error("CreatedAt or UpdatedAt not set after Create()")
	}
}

func TestGetByID(t *testing.T) {

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/avatar.jpg"

	expectedUser := &entities.User{

		ID:              userID,

		Name:            "Test User",

		Email:           "test@example.com",

		PasswordHash:    "hashedpassword",

		CreatedAt:       now,

		UpdatedAt:       now,

		IsEmailVerified: true,

		AvatarURL:       &avatarURL,

		IsDeleted:       false,

	}



	mockPool := &MockPgxPool{

		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {

			return &MockPgxRow{

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*uuid.UUID) = expectedUser.ID

					*dest[1].(*string) = expectedUser.Name

					*dest[2].(*string) = expectedUser.Email

					*dest[3].(*string) = expectedUser.PasswordHash

					*dest[4].(*time.Time) = expectedUser.CreatedAt

					*dest[5].(*time.Time) = expectedUser.UpdatedAt

					*dest[6].(*bool) = expectedUser.IsEmailVerified

					*dest[7].(*time.Time) = expectedUser.DeletedAt

					*dest[8].(*string) = *expectedUser.AvatarURL

					*dest[9].(*time.Time) = expectedUser.DeletionDueAt

					*dest[10].(*time.Time) = expectedUser.LastLoginAt

					*dest[11].(*bool) = expectedUser.IsDeleted

					return nil

				},

			}

		},

	}



	repo := NewUserRepository(mockPool)



	user, err := repo.GetByID(context.Background(), userID)

	if err != nil {

		t.Fatalf("GetByID() failed: %v", err)

	}



	if user.ID != expectedUser.ID {

		t.Errorf("Expected user ID %s, got %s", expectedUser.ID, user.ID)

	}

	// Add more assertions for other fields

}



func TestGetByEmail(t *testing.T) {

	userEmail := "test@example.com"

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/avatar.jpg"

	expectedUser := &entities.User{

		ID:              userID,

		Name:            "Test User",

		Email:           userEmail,

		PasswordHash:    "hashedpassword",

		CreatedAt:       now,

		UpdatedAt:       now,

		IsEmailVerified: true,

		AvatarURL:       &avatarURL,

		IsDeleted:       false,

	}



	mockPool := &MockPgxPool{

		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {

			return &MockPgxRow{

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*uuid.UUID) = expectedUser.ID

					*dest[1].(*string) = expectedUser.Name

					*dest[2].(*string) = expectedUser.Email

					*dest[3].(*string) = expectedUser.PasswordHash

					*dest[4].(*time.Time) = expectedUser.CreatedAt

					*dest[5].(*time.Time) = expectedUser.UpdatedAt

					*dest[6].(*bool) = expectedUser.IsEmailVerified

					*dest[7].(*time.Time) = expectedUser.DeletedAt

					*dest[8].(*string) = *expectedUser.AvatarURL

					*dest[9].(*time.Time) = expectedUser.DeletionDueAt

					*dest[10].(*time.Time) = expectedUser.LastLoginAt

					*dest[11].(*bool) = expectedUser.IsDeleted

					return nil

				},

			}

		},

	}



	repo := NewUserRepository(mockPool)



	user, err := repo.GetByEmail(context.Background(), userEmail)

	if err != nil {

		t.Fatalf("GetByEmail() failed: %t", err)

	}



	if user.Email != expectedUser.Email {

		t.Errorf("Expected user email %s, got %s", expectedUser.Email, user.Email)

	}

	// Add more assertions for other fields

}



func TestUpdateUser(t *testing.T) {

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/new_avatar.jpg"

	updatedUser := &entities.User{

		ID:              userID,

		Name:            "Updated Name",

		Email:           "updated@example.com",

		PasswordHash:    "newhashedpassword",

		IsEmailVerified: false,

		AvatarURL:       &avatarURL,

		DeletedAt:       nil,

		DeletionDueAt:   nil,

		LastLoginAt:     nil,

		IsDeleted:       false,

	}



	mockPool := &MockPgxPool{

		QueryRowFunc: func(ctx context.Context, sql string, args ...interface{}) pgx.Row {

			return &MockPgxRow{

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*time.Time) = now // updated_at

					return nil

				},

			}

		},

	}



	repo := NewUserRepository(mockPool)



	err := repo.Update(context.Background(), updatedUser)

	if err != nil {

		t.Fatalf("Update() failed: %v", err)

	}



	if updatedUser.UpdatedAt.IsZero() {

		t.Error("UpdatedAt not set after Update()")

	}

}



func TestDeleteUser(t *testing.T) {

	userID := uuid.New()



	mockPool := &MockPgxPool{

				ExecFunc: func(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {

					if len(arguments) > 0 && arguments[0].(uuid.UUID) == userID {

						return pgconn.CommandTag("UPDATE 1"), nil

					}

					return pgconn.CommandTag("UPDATE 0"), errors.New("user not found for deletion")

				},

			}



	repo := NewUserRepository(mockPool)



	err := repo.Delete(context.Background(), userID)

	if err != nil {

		t.Fatalf("Delete() failed: %v", err)

	}



	// Test case for user not found (or no rows affected)

	nonExistentUserID := uuid.New()

	err = repo.Delete(context.Background(), nonExistentUserID)

	if err == nil {

		t.Fatal("Expected an error for non-existent user, got nil")

	}

	if err.Error() != "user not found for deletion" {

		t.Errorf("Expected error 'user not found for deletion', got %v", err)

	}

}



func TestGetDeletedUsers(t *testing.T) {

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/avatar.jpg"

	expectedUser := &entities.User{

		ID:              userID,

		Name:            "Deleted User",

		Email:           "deleted@example.com",

		PasswordHash:    "hashedpassword",

		CreatedAt:       now,

		UpdatedAt:       now,

		IsEmailVerified: false,

		DeletedAt:       &now,

		AvatarURL:       &avatarURL,

		IsDeleted:       true,

	}



	mockPool := &MockPgxPool{

		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {

			return &MockPgxRows{

				NextFunc: func() bool { return true },

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*uuid.UUID) = expectedUser.ID

					*dest[1].(*string) = expectedUser.Name

					*dest[2].(*string) = expectedUser.Email

					*dest[3].(*string) = expectedUser.PasswordHash

					*dest[4].(*time.Time) = expectedUser.CreatedAt

					*dest[5].(*time.Time) = expectedUser.UpdatedAt

					*dest[6].(*bool) = expectedUser.IsEmailVerified

					*dest[7].(*time.Time) = expectedUser.DeletedAt

					*dest[8].(*string) = *expectedUser.AvatarURL

					*dest[9].(*time.Time) = expectedUser.DeletionDueAt

					*dest[10].(*time.Time) = expectedUser.LastLoginAt

					*dest[11].(*bool) = expectedUser.IsDeleted

					return nil

				},

				CloseFunc: func() {},

				ErrFunc:   func() error { return nil },

			}, nil

		},

	}



	repo := NewUserRepository(mockPool)

	users, err := repo.GetDeletedUsers(context.Background(), 10, 0)

	if err != nil {

		t.Fatalf("GetDeletedUsers() failed: %v", err)

	}

	if len(users) != 1 || users[0].ID != expectedUser.ID {

		t.Errorf("Expected 1 deleted user with ID %s, got %v", expectedUser.ID, users)

	}

}



func TestGetByEmailVerified(t *testing.T) {

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/avatar.jpg"

	expectedUser := &entities.User{

		ID:              userID,

		Name:            "Verified User",

		Email:           "verified@example.com",

		PasswordHash:    "hashedpassword",

		CreatedAt:       now,

		UpdatedAt:       now,

		IsEmailVerified: true,

		AvatarURL:       &avatarURL,

		IsDeleted:       false,

	}



	mockPool := &MockPgxPool{

		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {

			return &MockPgxRows{

				NextFunc: func() bool { return true },

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*uuid.UUID) = expectedUser.ID

					*dest[1].(*string) = expectedUser.Name

					*dest[2].(*string) = expectedUser.Email

					*dest[3].(*string) = expectedUser.PasswordHash

					*dest[4].(*time.Time) = expectedUser.CreatedAt

					*dest[5].(*time.Time) = expectedUser.UpdatedAt

					*dest[6].(*bool) = expectedUser.IsEmailVerified

					*dest[7].(*time.Time) = expectedUser.DeletedAt

					*dest[8].(*string) = *expectedUser.AvatarURL

					*dest[9].(*time.Time) = expectedUser.DeletionDueAt

					*dest[10].(*time.Time) = expectedUser.LastLoginAt

					*dest[11].(*bool) = expectedUser.IsDeleted

					return nil

				},

				CloseFunc: func() {},

				ErrFunc:   func() error { return nil },

			}, nil

		},

	}



	repo := NewUserRepository(mockPool)

	users, err := repo.GetByEmailVerified(context.Background(), true, 10, 0)

	if err != nil {

		t.Fatalf("GetByEmailVerified() failed: %v", err)

	}

	if len(users) != 1 || users[0].ID != expectedUser.ID {

		t.Errorf("Expected 1 verified user with ID %s, got %v", expectedUser.ID, users)

	}

}



func TestSearchByName(t *testing.T) {

	userID := uuid.New()

	now := time.Now()

	avatarURL := "http://example.com/avatar.jpg"

	expectedUser := &entities.User{

		ID:              userID,

		Name:            "Searchable User",

		Email:           "search@example.com",

		PasswordHash:    "hashedpassword",

		CreatedAt:       now,

		UpdatedAt:       now,

		IsEmailVerified: true,

		AvatarURL:       &avatarURL,

		IsDeleted:       false,

	}



	mockPool := &MockPgxPool{

		QueryFunc: func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {

			return &MockPgxRows{

				NextFunc: func() bool { return true },

				ScanFunc: func(dest ...interface{}) error {

					*dest[0].(*uuid.UUID) = expectedUser.ID

					*dest[1].(*string) = expectedUser.Name

					*dest[2].(*string) = expectedUser.Email

					*dest[3].(*string) = expectedUser.PasswordHash

					*dest[4].(*time.Time) = expectedUser.CreatedAt

					*dest[5].(*time.Time) = expectedUser.UpdatedAt

					*dest[6].(*bool) = expectedUser.IsEmailVerified

					*dest[7].(*time.Time) = expectedUser.DeletedAt

					*dest[8].(*string) = *expectedUser.AvatarURL

					*dest[9].(*time.Time) = expectedUser.DeletionDueAt

					*dest[10].(*time.Time) = expectedUser.LastLoginAt

					*dest[11].(*bool) = expectedUser.IsDeleted

					return nil

				},

				CloseFunc: func() {},

				ErrFunc:   func() error { return nil },

			}, nil

		},

	}



	repo := NewUserRepository(mockPool)

	users, err := repo.SearchByName(context.Background(), "Searchable", 10, 0)

	if err != nil {

		t.Fatalf("SearchByName() failed: %v", err)

	}

	if len(users) != 1 || users[0].ID != expectedUser.ID {

		t.Errorf("Expected 1 user with ID %s, got %v", expectedUser.ID, users)

	}

}
