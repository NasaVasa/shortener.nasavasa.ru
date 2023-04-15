package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"shortener/configs"
	r "shortener/router"
	"shortener/tests/mock"
	"testing"
)

func TestGetUrlsListHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Get urls list basic", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/get_urls_list/"+pgMock.Users[0].Token, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
	})

	t.Run("Get urls list wrong token", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/get_urls_list/wrong_token", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %v, got %v", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Get urls list wrong token time live", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		conf := configs.NewAuthConfig()
		pgMock.Users[0].TokenCreatedAt = pgMock.Users[0].TokenCreatedAt.Add(-(conf.TokenLiveTime + 1))
		req, err := http.NewRequest("GET", "/api/v1/get_urls_list/"+pgMock.Users[0].Token, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %v, got %v", http.StatusUnauthorized, w.Code)
		}
	})
}
