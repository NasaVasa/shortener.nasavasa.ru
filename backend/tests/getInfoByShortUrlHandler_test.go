package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"testing"
)

func TestGetInfoByShortUrlHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Get info basic", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/url_info/"+pgMock.Urls[0].ShortUrl, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
	})

	t.Run("Get info not found", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/url_info/wrongu", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusNotFound {
			t.Errorf("Expected status code %v, got %v", http.StatusNotFound, w.Code)
		}
	})
}
