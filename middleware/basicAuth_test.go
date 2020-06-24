package middleware

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
)

type user struct {
	Username string
	Password string
}

func validator(auth []byte) (bool, error) {
	var u user
	if err := json.Unmarshal(auth, &u); err != nil {
		return false, err
	}
	if u.Username == "test" && u.Password == "123123" {
		return true, nil
	}
	return false, nil

}

var testUser = map[string]string{"username": "test", "password": "123123"}

func TestBasicAuth(t *testing.T) {
	y := yee.New()
	y.Use(BasicAuth(validator))
	y.GET("/", func(context yee.Context) error {
		return context.String(http.StatusOK, "ok")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	y.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusUnauthorized, rec.Code)

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	u, _ := json.Marshal(testUser)
	encodeString := base64.StdEncoding.EncodeToString(u)
	req.Header.Set(yee.HeaderAuthorization, "basic "+encodeString)
	y.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusOK, rec.Code)

}
