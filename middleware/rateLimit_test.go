package middleware

import (
	"fmt"
	"github.com/cookieY/yee"
	"net/http"
	"net/http/httptest"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestRateLimit(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	y := yee.New()
	y.Use(RateLimitWithConfig(RateLimitConfig{Rate: 4,Time: time.Second * 2}))
	y.GET("/", func(context yee.Context) (err error) {
		return context.String(http.StatusOK, "ok")
	})
	var wg sync.WaitGroup
	var once sync.Once
	for i := 0; i < 10; i++ {
		if i > 5 {
			once.Do(func() {
				time.Sleep(time.Second * 1)
			})
		}
		wg.Add(1)
		go func(i int) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			y.ServeHTTP(rec, req)
			fmt.Printf("id: %d  code:%d \n",i,rec.Code)
			wg.Done()
		}(i)
	}
	wg.Wait()
}