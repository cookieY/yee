package middleware

import (
	"knocker"
)

type SecureConfig struct {
// Skipper defines a function to skip middleware.
Skipper Skipper

// XSSProtection provides protection against cross-site scripting attack (XSS)
// by setting the `X-XSS-Protection` header.
// Optional. Default value "1; mode=block".
XSSProtection string `yaml:"xss_protection"`

// ContentTypeNosniff provides protection against overriding Content-Type
// header by setting the `X-Content-Type-Options` header.
// Optional. Default value "nosniff".
ContentTypeNosniff string `yaml:"content_type_nosniff"`

// XFrameOptions can be used to indicate whether or not a browser should
// be allowed to render a page in a <frame>, <iframe> or <object> .
// Sites can use this to avoid clickjacking attacks, by ensuring that their
// content is not embedded into other sites.provides protection against
// clickjacking.
// Optional. Default value "SAMEORIGIN".
// Possible values:
// - "SAMEORIGIN" - The page can only be displayed in a frame on the same origin as the page itself.
// - "DENY" - The page cannot be displayed in a frame, regardless of the site attempting to do so.
// - "ALLOW-FROM uri" - The page can only be displayed in a frame on the specified origin.
XFrameOptions string `yaml:"x_frame_options"`

// HSTSMaxAge sets the `Strict-Transport-Security` header to indicate how
// long (in seconds) browsers should remember that this site is only to
// be accessed using HTTPS. This reduces your exposure to some SSL-stripping
// man-in-the-middle (MITM) attacks.
// Optional. Default value 0.
HSTSMaxAge int `yaml:"hsts_max_age"`

// HSTSExcludeSubdomains won't include subdomains tag in the `Strict Transport Security`
// header, excluding all subdomains from security policy. It has no effect
// unless HSTSMaxAge is set to a non-zero value.
// Optional. Default value false.
HSTSExcludeSubdomains bool `yaml:"hsts_exclude_subdomains"`

// ContentSecurityPolicy sets the `Content-Security-Policy` header providing
// security against cross-site scripting (XSS), clickjacking and other code
// injection attacks resulting from execution of malicious content in the
// trusted web page context.
// Optional. Default value "".
ContentSecurityPolicy string `yaml:"content_security_policy"`

// CSPReportOnly would use the `Content-Security-Policy-Report-Only` header instead
// of the `Content-Security-Policy` header. This allows iterative updates of the
// content security policy by only reporting the violations that would
// have occurred instead of blocking the resource.
// Optional. Default value false.
CSPReportOnly bool `yaml:"csp_report_only"`

// HSTSPreloadEnabled will add the preload tag in the `Strict Transport Security`
// header, which enables the domain to be included in the HSTS preload list
// maintained by Chrome (and used by Firefox and Safari): https://hstspreload.org/
// Optional.  Default value false.
HSTSPreloadEnabled bool `yaml:"hsts_preload_enabled"`

}

var (
	// DefaultSecureConfig is the default Secure middleware config.
	DefaultSecureConfig = SecureConfig{
		Skipper:            DefaultSkipper,
		XSSProtection:      "1; mode=block",
		ContentTypeNosniff: "nosniff",
		XFrameOptions:      "SAMEORIGIN",
		HSTSPreloadEnabled: false,
	}
)

func Xss() knocker.HandlerFunc {
	return func(c knocker.Context) error {
		c.SetHeader(knocker.HeaderXXSSProtection,"1; mode=block")
		return nil
	}
}
