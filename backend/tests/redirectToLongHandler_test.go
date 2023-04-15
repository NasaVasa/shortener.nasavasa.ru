package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"strings"
	"testing"
	"time"
)

func TestRedirectToLongHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Basic redirect to long link", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/"+pgMock.Urls[0].ShortUrl, nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code %v, got %v", 200, w.Code)
		}

		if pgMock.Urls[0].UrlClicks != 1 {
			t.Errorf("Expected 1 click, got %v", pgMock.Urls[0].UrlClicks)
		}
	})

	t.Run("Redirect to long link with wrong short url", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/wrong_short_url", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != 404 {
			t.Errorf("Expected status code %v, got %v", 404, w.Code)
		}
	})

	t.Run("Redirect to long link with time", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		wRed1 := httptest.NewRecorder()
		wRed2 := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(
				`{"url": "https://google.com","vip_key": "test", "time_to_live": 1, "time_to_live_unit": "SECONDS", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)

		req1, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(wRed1, req1)
		if wRed1.Code != 200 {
			t.Errorf("Expected status code %v, got %v", 200, w.Code)
		}

		for i, url := range pgMock.Urls {
			if url.ShortUrl == "test" {
				pgMock.Urls[i].UrlWillDelete = time.Now().UTC()
			}
		}

		req2, err := http.NewRequest("GET", "/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(wRed2, req2)
		if wRed2.Code != 404 {
			t.Errorf("Expected status code %v, got %v", 404, w.Code)
		}
	})
}
