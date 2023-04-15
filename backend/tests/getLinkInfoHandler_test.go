package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"testing"
	"time"
)

func TestLinkInfoHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Get not vip link info", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/admin/"+pgMock.Urls[0].SecretKey.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
	})

	t.Run("Get vip link info with time", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		wTime := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/api/v1/admin/"+pgMock.Urls[1].SecretKey.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}

		pgMock.Urls[1].UrlWillDelete = time.Now().UTC().Add(-time.Second)

		req, err = http.NewRequest("GET", "/api/v1/admin/"+pgMock.Urls[1].SecretKey.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(wTime, req)
		if wTime.Code != http.StatusNotFound {
			t.Errorf("Expected status code %v, got %v", http.StatusNotFound, wTime.Code)
		}
	})
}
