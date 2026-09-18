package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadAuth(t *testing.T) {
	e := echo.New()
	handler := downloadAuth("secret")(func(c echo.Context) error {
		return c.NoContent(http.StatusNoContent)
	})

	for _, test := range []struct {
		authorization string
		want          int
	}{
		{"", http.StatusUnauthorized},
		{"Bearer wrong", http.StatusUnauthorized},
		{"Bearer secret", http.StatusNoContent},
	} {
		req := httptest.NewRequest(http.MethodGet, "/dl/file.crt", nil)
		req.Header.Set(echo.HeaderAuthorization, test.authorization)
		rec := httptest.NewRecorder()
		err := handler(e.NewContext(req, rec))
		if err != nil {
			e.HTTPErrorHandler(err, e.NewContext(req, rec))
		}
		if rec.Code != test.want {
			t.Fatalf("Authorization %q returned %d, want %d", test.authorization, rec.Code, test.want)
		}
	}
}
