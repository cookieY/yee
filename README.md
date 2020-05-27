# Yee

![](https://img.shields.io/badge/build-alpha-brightgreen.svg)  
![](https://img.shields.io/badge/version-v0.0.1-brightgreen.svg)

Echo-like web framework. This is a framework for learning purposes. Study 7-golang and echo source code to realize some echo functions

-   Build RESTful APIs
-   Group APIs
-   Extensible middleware framework
-   Define middleware at root, group or route level
-   Data binding for JSON, XML and form payload

# Supported Go versions

Yee is available as a Go module. You need to use Go 1.1.3 +

## Example

```go
 	r := yee.New()

        r.GET("/", func(c yee.Context) error {
		return c.String(http.StatusOK, "<h1>Hello Gee</h1>")
	})

        r.Static("/assets", "dist/assets")

	r.GET("/", func(c yee.Context) (err error) {
		return c.HTMLTml(http.StatusOK, "dist/index.html")
	})

	r.POST("/test", func(c yee.Context) (err error) {
		u := new(p)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusOK, err.Error())
		}
		return c.JSON(http.StatusOK, u.Test)
	})
```

## License

MIT
