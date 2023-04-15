package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"strings"
	"testing"
)

func TestAuthHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Basic auth", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/auth", strings.NewReader(`{"login": "user1", "password": "11111"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
	})

	t.Run("Not found user", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/auth", strings.NewReader(`{"login": "user4", "password": "11111"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 404 {
			t.Errorf("Expected status code %v, got %v", 404, w.Code)
		}
	})

	t.Run("Wrong password", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/auth", strings.NewReader(`{"login": "user1", "password": "21111"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w.Code)
		}
	})
}
