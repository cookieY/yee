package middleware

//
//func TestLogger_LogWrite(t *testing.T) {
//
//	y := yee.New()
//	y.Logger.SetLevel(4)
//
//	y.POST("/hello/k/:b", func(c yee.Context) error {
//		y.Logger.Critical("有点问题")
//		y.Logger.Error("有点问题")
//		y.Logger.Warn("有点问题")
//		y.Logger.Info("有点问题")
//		y.Logger.Debug("有点问题")
//		return c.String(http.StatusOK, c.Params("b"))
//	})
//	t.Run("http_get", func(t *testing.T) {
//		req := httptest.NewRequest(http.MethodPost, "/hello/k/henry", nil)
//		rec := httptest.NewRecorder()
//		y.ServeHTTP(rec, req)
//		tx := assert.New(t)
//		tx.Equal("henry", rec.Body.String())
//		//assert.Equal("*", rec.Header().Get(yee.HeaderAccessControlAllowOrigin))
//	})
//}
