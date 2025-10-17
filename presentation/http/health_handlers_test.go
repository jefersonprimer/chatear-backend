package http_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jefersonprimer/chatear-backend/config"
	"github.com/jefersonprimer/chatear-backend/infrastructure"
	"github.com/jefersonprimer/chatear-backend/presentation/http"
	"github.com/stretchr/testify/assert"
)

func TestHealthz(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	infra, _ := infrastructure.NewInfrastructure("", "", "")
	cfg := &config.Config{}
	healthHandler := http.NewHealthHandler(infra, cfg)
	router.GET("/healthz", healthHandler.Healthz)

	req, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, `{"status":"UP"}`, w.Body.String())
}
