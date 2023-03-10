package middleware

import (
	"context"
	"fmt"
	"github.com/cookieY/yee"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

type ProxyTargetHandler interface {
	MatchAlgorithm() ProxyTarget
}

type ProxyTarget struct {
	Name string
	URL  *url.URL
}

type ProxyConfig struct {
	BalanceType    string
	Transport      http.RoundTripper
	ModifyResponse func(*http.Response) error
	ProxyTarget    ProxyTargetHandler
}

type errorHandler struct {
	Err  error
	Code int
}

func ProxyWithConfig(config ProxyConfig) yee.HandlerFunc {
	return func(c yee.Context) (err error) {
		req := c.Request()
		res := c.Response()
		if req.Header.Get(yee.HeaderXRealIP) == "" {
			req.Header.Set(yee.HeaderXRealIP, c.RemoteIP())
		}
		if req.Header.Get(yee.HeaderXForwardedProto) == "" {
			req.Header.Set(yee.HeaderXForwardedProto, c.Scheme())
		}
		if c.IsWebsocket() && req.Header.Get(yee.HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
			req.Header.Set(yee.HeaderXForwardedFor, c.RemoteIP())
		}
		switch {
		case c.IsWebsocket():
			//proxyRaw(tgt, c).ServeHTTP(res, req)
		case req.Header.Get(yee.HeaderAccept) == "text/event-stream":
		default:
			proxyHTTP(c, config).ServeHTTP(res, req)
		}
		if e, ok := c.Get("_error").(errorHandler); ok {
			return c.ServerError(e.Code, e.Err.Error())
		}

		return nil
	}
}

func proxyHTTP(c yee.Context, config ProxyConfig) http.Handler {
	tgt := config.ProxyTarget.MatchAlgorithm()
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		if err == context.Canceled || strings.Contains(err.Error(), "operation was canceled") {
			c.Put("_error", errorHandler{fmt.Errorf("client closed connection: %s", err.Error()), yee.StatusCodeContextCanceled})
		} else {
			c.Put("_error", errorHandler{fmt.Errorf("remote %s unreachable, could not forward: %v", desc, err), http.StatusBadGateway})
		}
	}
	proxy.Transport = config.Transport
	proxy.ModifyResponse = config.ModifyResponse
	return proxy
}
