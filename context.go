package yee

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cookieY/yee/logger"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const crashIndex int = math.MaxInt8 / 2

// Context is the default implementation  interface of context
type Context interface {
	Request() *http.Request
	Response() ResponseWriter
	HTML(code int, html string) (err error)
	JSON(code int, i interface{}) error
	ProtoBuf(code int, i proto.Message) error
	String(code int, s string) error
	FormValue(name string) string
	FormParams() (url.Values, error)
	FormFile(name string) (*multipart.FileHeader, error)
	File(file string) error
	Blob(code int, contentType string, b []byte) (err error)
	Status(code int)
	QueryParam(name string) string
	QueryString() string
	SetHeader(key string, value string)
	AddHeader(key string, value string)
	GetHeader(key string) string
	MultipartForm() (*multipart.Form, error)
	Redirect(code int, uri string) error
	Params(name string) string
	RequestURI() string
	Scheme() string
	IsTLS() bool
	IsWebsocket() bool
	Next()
	HTMLTpl(code int, tml string) (err error)
	QueryParams() map[string][]string
	Bind(i interface{}) error
	Cookie(name string) (*http.Cookie, error)
	SetCookie(cookie *http.Cookie)
	Cookies() []*http.Cookie
	Get(key string) interface{}
	Put(key string, values interface{})
	ServerError(code int, defaultMessage string) error
	RemoteIP() string
	Logger() logger.Logger
	Reset()
	Crash()
	IsCrash() bool
	CrashWithStatus(code int)
	CrashWithJson(code int, json interface{})
}

type context struct {
	engine    *Core
	writermem responseWriter
	w         ResponseWriter
	r         *http.Request
	path      string
	method    string
	code      int
	queryList url.Values // cache url.Values
	params    *Params
	Param     Params
	// middleware
	handlers  HandlersChain
	index     int
	store     map[string]interface{}
	lock      sync.RWMutex
	noRewrite bool
}

func (c *context) Reset() {
	c.index = -1
	c.handlers = c.engine.noRoute
}

func (c *context) Bind(i interface{}) error {
	return c.engine.bind.Bind(i, c)
}

func (c *context) reset() { // reset context members
	c.w = &c.writermem
	c.Param = c.Param[0:0]
	c.handlers = nil
	c.index = -1
	c.path = ""
	// when context reset clear queryList cache .
	// cause if not clear cache the queryParams results will mistake
	c.queryList = nil
	c.store = nil
	*c.params = (*c.params)[0:0]
}

func (c *context) Next() {
	c.index++
	s := len(c.handlers)
	if s > 0 {
		for ; c.index < s; c.index++ {
			_ = c.handlers[c.index](c)
			if c.w.Written() {
				break
			}
		}
	}
}

func (c *context) Logger() logger.Logger {
	return c.engine.l
}

func (c *context) ServerError(code int, defaultMessage string) error {
	c.writermem.status = code
	if c.writermem.Written() {
		return errors.New("headers were already written")
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = []string{MIMETextPlainCharsetUTF8}
		c.Logger().Error(fmt.Sprintf("%s %s", c.r.URL, defaultMessage))
		_, err := c.w.Write([]byte(defaultMessage))
		if err != nil {
			return fmt.Errorf("cannot write message to writer during serve error: %v", err)
		}
		return nil
	}
	c.writermem.WriteHeaderNow()
	return nil
}

func (c *context) Put(key string, values interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.store == nil {
		c.store = make(map[string]interface{})
	}
	c.store[key] = values
}

func (c *context) Crash() {
	c.index = crashIndex
}

func (c *context) IsCrash() bool {
	return c.index >= crashIndex
}

func (c *context) CrashWithStatus(code int) {
	c.Status(code)
	c.w.WriteHeaderNow()
	c.Crash()
}

func (c *context) CrashWithJson(code int, json interface{}) {
	c.Crash()
	err := c.JSON(code, json)
	if err != nil {
		return
	}
}

func (c *context) Get(key string) interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *context) Request() *http.Request {
	return c.r
}

func (c *context) Response() ResponseWriter {
	return c.w
}

func (c *context) RemoteIP() string {
	if ip := c.r.Header.Get(HeaderXForwardedFor); ip != "" {
		i := strings.IndexAny(ip, ", ")
		if i > 0 {
			return ip[:i]
		}
	}
	if ip := c.r.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ip, _, _ := net.SplitHostPort(c.r.RemoteAddr)
	return ip
}

func (c *context) HTML(code int, html string) (err error) {
	return c.HTMLBlob(code, []byte(html))
}

func (c *context) HTMLTpl(code int, tml string) (err error) {
	s, e := ioutil.ReadFile(tml)
	if e != nil {
		panic(e)
	}
	return c.HTMLBlob(code, s)
}

func (c *context) HTMLBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMETextHTMLCharsetUTF8, b)
}

func (c *context) Blob(code int, contentType string, b []byte) (err error) {
	if !c.writermem.Written() {
		c.writeContentType(contentType)
		c.w.WriteHeader(code)
		if _, err = c.w.Write(b); err != nil {
			c.Logger().Error(err.Error())
		}
	}
	return
}

func (c *context) JSON(code int, i interface{}) (err error) {
	if !c.writermem.Written() {
		enc := json.NewEncoder(c.w)
		c.writeContentType(MIMEApplicationJSONCharsetUTF8)
		c.w.WriteHeader(code)
		return enc.Encode(i)
	}
	return
}

func (c *context) ProtoBuf(code int, i proto.Message) (err error) {
	if !c.writermem.Written() {
		c.writeContentType(MIMEApplicationProtobuf)
		c.w.WriteHeader(code)
		b, err := proto.Marshal(i)
		if err != nil {
			return err
		}
		if _, err = c.w.Write(b); err != nil {
			return err
		}
	}
	return
}

func (c *context) JSONP(code int, fn string, i interface{}) error {
	enc := json.NewEncoder(c.w)
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.w.WriteHeader(code)
	if _, err := c.w.Write([]byte(fn + "(")); err != nil {
		return err
	}
	if err := enc.Encode(i); err != nil {
		return err
	}
	if _, err := c.w.Write([]byte(");")); err != nil {
		return err
	}
	return nil
}

func (c *context) String(code int, s string) error {
	return c.Blob(code, MIMETextPlainCharsetUTF8, StringToBytes(s))
}

func (c *context) Status(code int) {
	c.w.WriteHeader(code)
}

func (c *context) SetHeader(key string, value string) {
	c.w.Header().Set(key, value)
}

func (c *context) AddHeader(key string, value string) {
	c.w.Header().Add(key, value)
}

func (c *context) GetHeader(key string) string {
	return c.r.Header.Get(key)
}

func (c *context) Params(name string) string {
	for _, i := range *c.params {
		if i.Key == name {
			return i.Value
		}
	}
	return ""
}
func (c *context) QueryParams() map[string][]string {
	return c.r.URL.Query()
}

func (c *context) QueryParam(name string) string {
	if c.queryList == nil {
		c.queryList = c.r.URL.Query()
	}
	return c.queryList.Get(name)
}

func (c *context) QueryString() string {
	return c.r.URL.RawQuery
}

func (c *context) FormValue(name string) string {
	return c.r.FormValue(name)
}

func (c *context) FormParams() (url.Values, error) {
	if strings.HasPrefix(c.r.Header.Get(HeaderContentType), MIMEMultipartForm) {
		if err := c.r.ParseMultipartForm(defaultMemory); err != nil {
			return nil, err
		}
	} else {
		if err := c.r.ParseForm(); err != nil {
			return nil, err
		}
	}
	return c.r.Form, nil
}

func (c *context) FormFile(name string) (*multipart.FileHeader, error) {
	_, fd, err := c.r.FormFile(name)
	return fd, err
}

func (c *context) File(file string) error {
	fd, err := os.Open(file)
	if err != nil {
		return err
	}

	defer fd.Close()

	f, _ := fd.Stat()
	if f.IsDir() {
		file = filepath.Join(file, indexPage)
		fd, err = os.Open(file)
		if err != nil {
			return ErrNotFoundHandler
		}
		defer fd.Close()
		if f, err = fd.Stat(); err != nil {
			return err
		}
	}
	http.ServeContent(c.Response(), c.Request(), f.Name(), f.ModTime(), fd)
	return nil
}

func (c *context) MultipartForm() (*multipart.Form, error) {
	err := c.r.ParseMultipartForm(defaultMemory)
	return c.r.MultipartForm, err
}

func (c *context) Cookie(name string) (*http.Cookie, error) {
	return c.r.Cookie(name)
}

func (c *context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.w, cookie)
}

func (c *context) Cookies() []*http.Cookie {
	return c.r.Cookies()
}

func (c *context) RequestURI() string {
	return c.r.RequestURI
}

func (c *context) Scheme() string {
	scheme := "http"
	if scheme := c.r.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.r.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.r.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.r.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return scheme
}

func (c *context) IsTLS() bool {
	return c.r.TLS != nil
}

func (c *context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.r.Header.Get(HeaderConnection)), "upgrade") &&
		strings.EqualFold(c.r.Header.Get(HeaderUpgrade), "websocket") {
		return true
	}
	return false
}

func (c *context) Redirect(code int, uri string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.w.Header().Set(HeaderLocation, uri)
	c.w.WriteHeader(code)
	return nil
}

func (c *context) writeContentType(value string) {
	header := c.w.Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}
