package main

import (
	"fmt"
	"net/http"

	"github.com/Dionizio8/go-rate-limiter/configs"
	"github.com/Dionizio8/go-rate-limiter/pkg/rtl"
	"github.com/Dionizio8/go-rate-limiter/pkg/rtl/infra/database"
	middlewarertl "github.com/Dionizio8/go-rate-limiter/pkg/rtl/middleware"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/redis/go-redis/v9"
)

func main() {
	cfg, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprint(cfg.REDIS_HOST, ":", cfg.REDIS_PORT),
		DB:   0,
	})

	rtlRedisRepository := database.NewRateLimitRedisRepository(rdb)
	rtl := rtl.NewRTL(rtlRedisRepository, cfg.RTLIP, cfg.RTLBlockTime)
	for _, token := range cfg.RTLTokens {
		rtl.SetToken(token.Token, token.ExpirationTime)
	}

	middlewareRTL := middlewarertl.NewMiddlewareRTL(rtl)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middlewareRTL.RateLimit)
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello GoExpert 🚀"))
	})

	http.ListenAndServe(cfg.WebServerPort, r)
}
