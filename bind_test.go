package yee

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type user struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password"`
	Age      int    `json:"age"`
}

type empty struct {
}

type cmdbBind struct {
	RegionId string `json:"region_id"`
	SecId    string `json:"secId"`
	Cloud    string `json:"cloud"`
	Account  string `json:"account"`
}

var userInfo = `{"username": "henry","age":24,"password":"123123"}`
var invalidInfo = `{"username": "","age":24,"password":"123123"}`
var encrypt = `e2db79dc56e0b5a5866fa4062c9c715e66a5d4820d5424c7645092be0041c1e62d9571b0549758cd02445593b2a276d455cba5b31295e1288d67255e78e4dd78`

func TestBindJSON(t *testing.T) {
	assertions := assert.New(t)
	testBindOkay(assertions, strings.NewReader(encrypt), MIMEApplicationJSON)
	//testBindError(assertions, strings.NewReader(invalidInfo), MIMEApplicationJSON)
	//testBindQueryPrams(assertions, MIMETextHTML)
}

func TestDefaultBinder_Bind(t *testing.T) {
	e := C()
	e.POST("/bind", func(c Context) (err error) {
		u := new(user)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	req := httptest.NewRequest(http.MethodPost, "/bind", strings.NewReader(invalidInfo))
	req.Header.Set("Content-Type", MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	fmt.Println(rec.Code)
	fmt.Println(rec.Body.String())
}

func TestBindEncryptOkay(t *testing.T) {
	e := C()
	e.POST("/", func(c Context) (err error) {
		u := new(user)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, u)
	})
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(encrypt))
	req.Header.Set("Content-Type", MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
}

func TestDefaultBinder_Params_Bind(t *testing.T) {
	e := C()
	e.GET("/bind", func(c Context) (err error) {
		u := new(empty)
		if err := c.Bind(u); err != nil {
			return err
		}
		return c.JSON(http.StatusOK, "")
	})
	req := httptest.NewRequest(http.MethodGet, "/bind?username=xxxx&cn", nil)
	req.Header.Set("Content-Type", MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func testBindOkay(assert *assert.Assertions, r io.Reader, ctype string) {
	e := C()
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
	e := C()
	req := httptest.NewRequest(http.MethodPost, "/", r)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set(HeaderContentType, ctype)
	u := new(user)
	err := c.Bind(u)
	assert.Error(err, "Unmarshal type error: expected=yee.user, got=number, field=, offset=1")
}

func testBindQueryPrams(assert *assert.Assertions, ctype string) {
	e := C()
	req := httptest.NewRequest(http.MethodGet, "/?secId=sg-gw86k2rjop30v1ktyn3j&region_id=eu-central", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	req.Header.Set(HeaderContentType, ctype)
	u := new(cmdbBind)
	err := c.Bind(u)
	if assert.NoError(err) {
		assert.Equal("eu-central", u.RegionId)
		assert.Equal("sg-gw86k2rjop30v1ktyn3j", u.SecId)
	}
}

func TestQueryParams(t *testing.T) {
	assertions := assert.New(t)
	testBindQueryPrams(assertions, MIMETextHTML)
}
