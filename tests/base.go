package tests

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/vesicash/mor-api/internal/config"
	"github.com/vesicash/mor-api/internal/models/migrations"
	"github.com/vesicash/mor-api/pkg/repository/storage/postgresql"
	"github.com/vesicash/mor-api/utility"
)

func Setup() *utility.Logger {
	logger := utility.NewLogger()
	config := config.Setup(logger, "../../app")
	db := postgresql.ConnectToDatabases(logger, config.TestDatabases)
	if config.TestDatabases.Migrate {
		migrations.RunAllMigrations(db)
	}
	return logger
}

func ParseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func AssertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("handler returned wrong status code: got status %d expected status %d", got, expected)
	}
}

func AssertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("handler returned wrong message: got message: %q expected: %q", got, expected)
	}
}
