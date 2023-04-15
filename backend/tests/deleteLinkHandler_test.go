package tests

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	r "shortener/router"
	"shortener/tests/mock"
	"testing"
)

func TestDeleteLinkHandler(t *testing.T) {
	pgMock := mock.NewPgMock()
	router := gin.Default()
	r.SetupRouter(router, pgMock)

	t.Run("Delete link basic", func(t *testing.T) {
		defer mock.ClearMock(pgMock)
		mock.FillTestDataInMock(pgMock)
		basicLen := len(pgMock.Urls)
		w := httptest.NewRecorder()
		req, err := http.NewRequest("DELETE", "/api/v1/admin/"+pgMock.Urls[0].SecretKey.String(), nil)
		if err != nil {
			t.Fatal(err)
		}
		router.ServeHTTP(w, req)
		if w.Code != 204 {
			t.Errorf("Expected status code %v, got %v", 204, w.Code)
		}
		if len(pgMock.Urls) != basicLen-1 {
			t.Errorf("Expected length of urls %v, got %v", basicLen-1, len(pgMock.Urls))
		}
	})
}
