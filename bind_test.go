package yee

import (
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

var userInfo = `{"username": "henry","age":24,"password":"123123"}`
var invalidInfo = `1{"username": "henry","age":24,"password":"123123"}`
func TestBindJSON(t *testing.T) {
	assertions := assert.New(t)
	testBindOkay(assertions, strings.NewReader(userInfo), MIMEApplicationJSON)
	testBindError(assertions,strings.NewReader(invalidInfo),MIMEApplicationJSON)
}

func testBindOkay(assert *assert.Assertions, r io.Reader, ctype string) {
	e := New()
	req := httptest.NewRequest(http.MethodPost, "/", r)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set(HeaderContentType, ctype)
	u := new(user)
	err := c.Bind(u)
	if assert.NoError(err) {
		assert.Equal("henry", u.Username)
		assert.Equal(24, u.Age)
		assert.Equal("123123", u.Password)
	}
}

func testBindError(assert *assert.Assertions, r io.Reader, ctype string) {
	e := New()
	req := httptest.NewRequest(http.MethodPost, "/", r)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set(HeaderContentType, ctype)
	u := new(user)
	err := c.Bind(u)
		switch {
		case strings.HasPrefix(ctype, MIMEApplicationJSON), strings.HasPrefix(ctype, MIMEApplicationXML), strings.HasPrefix(ctype, MIMETextXML),
			strings.HasPrefix(ctype, MIMEApplicationForm), strings.HasPrefix(ctype, MIMEMultipartForm):
				assert.Equal(http.StatusBadRequest, rec.Code)
		default:
			assert.Equal(ErrUnsupportedMediaType, err)
		}
}
