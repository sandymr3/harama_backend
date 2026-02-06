package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"harama/internal/api"
	"harama/internal/config"
	"harama/internal/repository/postgres"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{Port: "8080"}
	db, _ := postgres.Connect("postgres://test")
	router, _ := api.NewRouter(cfg, db)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d", w.Code)
	}
}

func TestCreateExam(t *testing.T) {
	t.Skip("Integration test - requires DB")
}

func TestGradingEngine(t *testing.T) {
	t.Skip("Integration test - requires Gemini API")
}
