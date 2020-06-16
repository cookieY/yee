package middleware

import (
	"encoding/base64"
	"encoding/json"
	"github.com/cookieY/yee"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

type user struct {
	Username string
	Password string
}

func validator(auth []byte) (error, bool) {
	var u user
	if err := json.Unmarshal(auth, &u); err != nil {
		return err, false
	}
	if u.Username == "test" && u.Password == "123123" {
		return nil, true
	}
	return nil, false

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

