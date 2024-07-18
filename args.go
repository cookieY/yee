package yee

import (
	"errors"
)

// Header types
const (
	HeaderSecWebSocketProtocol = "Sec-Websocket-Protocol"
	HeaderAccept               = "Accept"
	HeaderAcceptEncoding       = "Accept-Encoding"
	HeaderAuthorization        = "Authorization"
	HeaderContentDisposition   = "Content-Disposition"
	HeaderContentEncoding      = "Content-Encoding"
	HeaderContentLength        = "Content-Length"
	HeaderContentType          = "Content-Type"
	HeaderCookie               = "Cookie"
	HeaderSetCookie            = "Set-Cookie"
	HeaderIfModifiedSince      = "If-Modified-Since"
	HeaderLastModified         = "Last-Modified"
	HeaderLocation             = "Location"
	HeaderUpgrade              = "Upgrade"
	HeaderConnection           = "Connection"
	HeaderVary                 = "Vary"
	HeaderWWWAuthenticate      = "WWW-Authenticate"
	HeaderXForwardedFor        = "X-Forwarded-For"
	HeaderXForwardedProto      = "X-Forwarded-Proto"
	HeaderXForwardedProtocol   = "X-Forwarded-Protocol"
	HeaderXForwardedSsl        = "X-Forwarded-Ssl"
	HeaderXUrlScheme           = "X-Url-Scheme"
	HeaderXHTTPMethodOverride  = "X-HTTP-Method-Override"
	HeaderXRealIP              = "X-Real-IP"
	HeaderXRequestID           = "X-Request-ID"
	HeaderXRequestedWith       = "X-Requested-With"
	HeaderServer               = "Server"
	HeaderOrigin               = "Origin"

	// Access control
	HeaderAccessControlRequestMethod    = "Access-Control-Request-Method"
	HeaderAccessControlRequestHeaders   = "Access-Control-Request-Headers"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
	HeaderAccessControlMaxAge           = "Access-Control-Max-Age"

	// Security
	HeaderStrictTransportSecurity         = "Strict-Transport-Security"
	HeaderXContentTypeOptions             = "X-Content-Type-Options"
	HeaderXXSSProtection                  = "X-XSS-Protection"
	HeaderXFrameOptions                   = "X-Frame-Options"
	HeaderContentSecurityPolicy           = "Content-Security-Policy"
	HeaderContentSecurityPolicyReportOnly = "Content-Security-Policy-Report-Only"
	HeaderXCSRFToken                      = "X-CSRF-Token"
	HeaderReferrerPolicy                  = "Referrer-Policy"
)

const (
	defaultMemory = 32 << 20 // 32 MB
	indexPage     = "index.html"
	defaultIndent = "  "
)

const (
	StatusCodeContextCanceled = 499
)

// MIME types
const (
	MIMEApplicationJSON                  = "application/json"
	MIMEApplicationJSONCharsetUTF8       = MIMEApplicationJSON + "; " + charsetUTF8
	MIMEApplicationJavaScript            = "application/javascript"
	MIMEApplicationJavaScriptCharsetUTF8 = MIMEApplicationJavaScript + "; " + charsetUTF8
	MIMEApplicationXML                   = "application/xml"
	MIMEApplicationXMLCharsetUTF8        = MIMEApplicationXML + "; " + charsetUTF8
	MIMETextXML                          = "text/xml"
	MIMETextXMLCharsetUTF8               = MIMETextXML + "; " + charsetUTF8
	MIMEApplicationForm                  = "application/x-www-form-urlencoded"
	MIMEApplicationProtobuf              = "application/protobuf"
	MIMEApplicationMsgpack               = "application/msgpack"
	MIMETextHTML                         = "text/html"
	MIMETextHTMLCharsetUTF8              = MIMETextHTML + "; " + charsetUTF8
	MIMETextPlain                        = "text/plain"
	MIMETextPlainCharsetUTF8             = MIMETextPlain + "; " + charsetUTF8
	MIMEMultipartForm                    = "multipart/form-data"
	MIMEOctetStream                      = "application/octet-stream"
)

const (
	charsetUTF8 = "charset=UTF-8"
	serverName  = "yee"
)

// Err types
var (
	ErrUnsupportedMediaType   = errors.New("http server not support media type")
	ErrValidatorNotRegistered = errors.New("validator not registered")
	ErrRendererNotRegistered  = errors.New("renderer not registered")
	ErrInvalidRedirectCode    = errors.New("invalid redirect status code")
	ErrCookieNotFound         = errors.New("cookie not found")
	ErrNotFoundHandler        = errors.New("404 NOT FOUND")
	ErrInvalidCertOrKeyType   = errors.New("invalid cert or key type, must be string or []byte")
)
