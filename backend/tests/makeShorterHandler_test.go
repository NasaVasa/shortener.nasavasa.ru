package tests

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"strings"
	"testing"
)

func TestMakeShorterHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Basic create short link", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter", strings.NewReader(`{"url": "https://google.com"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		if len(pgMock.Urls) != 1 {
			t.Errorf("Expected 1 url, got %v", len(pgMock.Urls))
		}
	})

	t.Run("Basic create short link without http", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter", strings.NewReader(`{"url": "a"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		if len(pgMock.Urls) != 1 {
			t.Errorf("Expected 1 url, got %v", len(pgMock.Urls))
		}
	})

	t.Run("Create short links with the same names", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		w := httptest.NewRecorder()
		wAdd := httptest.NewRecorder()
		req1, err := http.NewRequest("POST", "/api/v1/make_shorter", strings.NewReader(`{"url": "https://google.com"}`))
		if err != nil {
			t.Fatal(err)
		}
		req1.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req1)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		res1 := w.Body.String()

		req2, err := http.NewRequest("POST", "/api/v1/make_shorter", strings.NewReader(`{"url": "https://google.com"}`))
		if err != nil {
			t.Fatal(err)
		}

		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wAdd, req2)
		if wAdd.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		res2 := wAdd.Body.String()
		if len(pgMock.Urls) != 1 {
			t.Errorf("Expected 1 url, got %v", len(pgMock.Urls))
		}

		if res1 != res2 {
			t.Errorf("Expected the same short url, got %v and %v", res1, res2)
		}
	})

	t.Run("Create Vip link basic test", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		baseLength := len(pgMock.Urls)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(
				`{"url": "https://google.com","vip_key": "test", "time_to_live": 1000, "time_to_live_unit": "SECONDS", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %v, got %v", http.StatusOK, w.Code)
		}
		if len(pgMock.Urls) != baseLength+1 {
			t.Errorf("Expected %v url, got %v", baseLength+1, len(pgMock.Urls))
		}
	})

	t.Run("Create Vip link with wrong token", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		randomToken := uuid.NewString()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://google.com","vip_key": "test", "time_to_live": 1000, "time_to_live_unit": "SECONDS", "token": "`+randomToken+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != http.StatusUnauthorized {
			t.Errorf("Expected status code %v, got %v", http.StatusUnauthorized, w.Code)
		}
	})

	t.Run("Create Vip link with existing vip key", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://google.com","vip_key": "yandex", "time_to_live": 1000, "time_to_live_unit": "SECONDS", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w.Code)
		}
	})

	t.Run("Create Vip link with wrong time to live and time to live unit", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		baseLen := len(pgMock.Urls)
		w := httptest.NewRecorder()
		wAdd := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://google.com","vip_key": "test", "time_to_live": 1000, "time_to_live_unit": "WRONG", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w.Code)
		}

		req2, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://google.com","vip_key": "test", "time_to_live": -2, "time_to_live_unit": "SECONDS", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wAdd, req2)
		if w.Code != 400 {
			t.Errorf("Expected status code %v, got %v", 400, w.Code)
		}

		if len(pgMock.Urls) != baseLen {
			t.Errorf("Expected %v url, got %v", baseLen, len(pgMock.Urls))
		}
	})

	t.Run("Create generating short link and vip link with the same long url", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		w := httptest.NewRecorder()
		wAdd := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://test.com"}`))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(w, req)
		if w.Code != 200 {
			t.Errorf("Expected status code %v, got %v", 200, w.Code)
		}

		req2, err := http.NewRequest("POST", "/api/v1/make_shorter",
			strings.NewReader(`{"url": "https://test.com","vip_key": "test", "time_to_live": 20, "time_to_live_unit": "SECONDS", "token": "`+pgMock.Users[0].Token+`"}`))
		if err != nil {
			t.Fatal(err)
		}
		req2.Header.Set("Content-Type", "application/json")
		router.ServeHTTP(wAdd, req2)
		if w.Code != 200 {
			t.Errorf("Expected status code %v, got %v", 200, w.Code)
		}
	})
}
