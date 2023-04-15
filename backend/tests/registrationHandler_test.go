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

func TestRegistrationHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Basic registration", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/registration", strings.NewReader(`{"login": "test", "email":"test@gmail.com", "password": "test"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		if len(pgMock.Users) != 1 {
			t.Errorf("Expected 1 user, got %v", len(pgMock.Users))
		}
	})

	t.Run("Check wrong email format", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/registration", strings.NewReader(`{"login": "test", "email":"testgmail.com", "password": "test"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w.Code)
		}
		if len(pgMock.Users) != 0 {
			t.Errorf("Expected 0 user, got %v", len(pgMock.Users))
		}
	})

	t.Run("User already exist by login test", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w1 := httptest.NewRecorder()
		w2 := httptest.NewRecorder()
		mock.FillTestDataInMock(pgMock)
		req1, err := http.NewRequest("POST", "/api/v1/registration", strings.NewReader(`{"login": "user1", "email":"test@gmail.com",password: "test"}`))
		if err != nil {
			t.Fatal(err)
		}
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w1, req1)
		if w1.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w1.Code)
		}

		req2, err := http.NewRequest("POST", "/api/v1/registration", strings.NewReader(`{"login": "usertest", "email":"example1@gmail.com", "password": "test"}`))
		if err != nil {
			t.Fatal(err)
		}
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w2, req2)
		if w2.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w2.Code)
		}
	})

}
