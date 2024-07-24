package middleware

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/Dionizio8/go-rate-limiter/pkg/rtl"
	"github.com/Dionizio8/go-rate-limiter/pkg/rtl/infra/database"
	"github.com/go-chi/chi"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

const redis_addr = "localhost:6379"

var rdb *redis.Client

func TestMain(m *testing.M) {
	rdb = redis.NewClient(&redis.Options{
		Addr: redis_addr,
		DB:   0,
	})

	code := m.Run()
	os.Exit(code)
}

func TestMiddleware_RateLimit_IPBlock(t *testing.T) {
	rtlRedisRepository := database.NewRateLimitRedisRepository(rdb)
	rtl := rtl.NewRTL(rtlRedisRepository, 1, 3)
	middlewareRTL := NewMiddlewareRTL(rtl)

	r := chi.NewRouter()
	r.Use(middlewareRTL.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test 123..."))
	})

	// Request 1
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "test 123...", w1.Body.String())

	// Request 2 (should be blocked)
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Equal(t, errorRT, w2.Body.String())

	time.Sleep(3 * time.Second)

	// Request 3 (should be allowed)
	req3 := httptest.NewRequest(http.MethodGet, "/", nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	assert.Equal(t, "test 123...", w3.Body.String())
}

func TestMiddleware_RateLimit_TokenBlock(t *testing.T) {
	rtlRedisRepository := database.NewRateLimitRedisRepository(rdb)
	rtl := rtl.NewRTL(rtlRedisRepository, 1, 2)
	rtl.SetToken("token1", 4)
	middlewareRTL := NewMiddlewareRTL(rtl)

	r := chi.NewRouter()
	r.Use(middlewareRTL.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test 123..."))
	})

	// Request 1
	req1 := httptest.NewRequest(http.MethodGet, "/", nil)
	req1.Header.Set("API_KEY", "token1")
	w1 := httptest.NewRecorder()
	r.ServeHTTP(w1, req1)

	assert.Equal(t, http.StatusOK, w1.Code)
	assert.Equal(t, "test 123...", w1.Body.String())

	// Request 2
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	req2.Header.Set("API_KEY", "token1")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusTooManyRequests, w2.Code)
	assert.Equal(t, errorRT, w2.Body.String())

	time.Sleep(2 * time.Second)

	// Request 3
	req3 := httptest.NewRequest(http.MethodGet, "/", nil)
	req3.Header.Set("API_KEY", "token1")
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusTooManyRequests, w3.Code)
	assert.Equal(t, errorRT, w3.Body.String())

	time.Sleep(3 * time.Second)

	// Request 4
	req4 := httptest.NewRequest(http.MethodGet, "/", nil)
	req4.Header.Set("API_KEY", "token1")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	assert.Equal(t, http.StatusOK, w4.Code)
	assert.Equal(t, "test 123...", w4.Body.String())
}
