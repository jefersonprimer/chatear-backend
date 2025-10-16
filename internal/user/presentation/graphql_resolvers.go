package presentation

import (
	"context"

	"github.com/google/uuid"
	"github.com/jefersonprimer/chatear-backend/internal/user/application"
	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
)

// Resolver is the root resolver for the GraphQL schema.
type Resolver struct {
	UserHandler *UserHandler
}

// NewResolver creates a new Resolver.
func NewResolver(userHandler *UserHandler) *Resolver {
	return &Resolver{UserHandler: userHandler}
}

// Mutation returns the mutation resolver.
func (r *Resolver) Mutation() MutationResolver {
	return &mutationResolver{r}
}

// Query returns the query resolver.
func (r *Resolver) Query() QueryResolver {
	return &queryResolver{r}
}

// mutationResolver is the resolver for the Mutation type.
type mutationResolver struct{ *Resolver }

// Register resolves the register mutation.
func (r *mutationResolver) Register(ctx context.Context, input struct{ Name, Email, Password string }) (*domain.User, error) {
	return r.UserHandler.RegisterUser.Execute(ctx, input.Name, input.Email, input.Password)
}

// Login resolves the login mutation.
func (r *mutationResolver) Login(ctx context.Context, input struct{ Email, Password string }) (*application.LoginResponse, error) {
	return r.UserHandler.LoginUser.Execute(ctx, input.Email, input.Password)
}

// Logout resolves the logout mutation.
func (r *mutationResolver) Logout(ctx context.Context, input struct{ Token string }) (*bool, error) {
	err := r.UserHandler.LogoutUser.Execute(ctx, input.Token)
	if err != nil {
		return nil, err
	}
	result := true
	return &result, nil
}

// DeleteUser resolves the deleteUser mutation.
func (r *mutationResolver) DeleteUser(ctx context.Context, id string) (*bool, error) {
	userID, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	err = r.UserHandler.DeleteUser.Execute(ctx, userID)
	if err != nil {
		return nil, err
	}
	result := true
	return &result, nil
}

// queryResolver is the resolver for the Query type.
type queryResolver struct{ *Resolver }

// Me resolves the me query.
func (r *queryResolver) Me(ctx context.Context) (*domain.User, error) {
	// In a real application, you would get the user ID from the context
	// (e.g., from a JWT token) and then retrieve the user from the database.
	// For now, we'll just return a placeholder user.
	return &domain.User{ID: uuid.New(), Name: "Placeholder", Email: "placeholder@example.com"}, nil
}

// MutationResolver is the interface for the Mutation type.
type MutationResolver interface {
	Register(ctx context.Context, input struct{ Name, Email, Password string }) (*domain.User, error)
	Login(ctx context.Context, input struct{ Email, Password string }) (*application.LoginResponse, error)
	Logout(ctx context.Context, input struct{ Token string }) (*bool, error)
	DeleteUser(ctx context.Context, id string) (*bool, error)
}

// QueryResolver is the interface for the Query type.
type QueryResolver interface {
	Me(ctx context.Context) (*domain.User, error)
}
