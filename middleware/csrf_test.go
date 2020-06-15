package middleware

import (
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCSRFWithConfig(t *testing.T) {

	y := yee.InitCore()
	y.Use(Csrf())
	y.POST("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "ok")
	})
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rec := httptest.NewRecorder()

	// Without CSRF cookie
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "missing csrf token in the header string", rec.Body.String())
	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// invalid csrf token
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(yee.HeaderXCSRFToken, "cbghjiwhd")
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, "invalid csrf token", rec.Body.String())
	assert.Equal(t, http.StatusForbidden, rec.Code)

	token := yee.RandomString(16)
	req = httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(yee.HeaderCookie, "_csrf="+token)
	req.Header.Set(yee.HeaderXCSRFToken, token)
	rec = httptest.NewRecorder()
	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)
}