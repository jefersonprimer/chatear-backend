package constants

// Application constants
const (
	// API versions
	APIVersion = "v1"
	
	// HTTP status codes
	StatusOK                  = 200
	StatusCreated             = 201
	StatusBadRequest          = 400
	StatusUnauthorized        = 401
	StatusForbidden           = 403
	StatusNotFound            = 404
	StatusConflict            = 409
	StatusInternalServerError = 500
	
	// Database
	DefaultPageSize = 20
	MaxPageSize     = 100
	
	// JWT
	JWTHeaderName = "Authorization"
	JWTBearer           = "Bearer"
	JwtSecret           = "your-secret-key" // TODO: Replace with a strong, environment-variable-based secret
	AccessTokenExpiration = 15 * time.Minute
	
	// Cache
	DefaultCacheTTL = 3600 // 1 hour in seconds
)