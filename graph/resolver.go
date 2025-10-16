package graph

import (
	"github.com/jefersonprimer/chatear-backend/internal/user/application"
	"github.com/jefersonprimer/chatear-backend/shared/auth"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	UserAppService *application.UserApplicationService
	TokenService *auth.TokenService
}

