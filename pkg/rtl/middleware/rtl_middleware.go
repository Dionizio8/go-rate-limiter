package middleware

import (
	"net/http"

	"github.com/Dionizio8/go-rate-limiter/pkg/rtl"
)

const (
	errorRT = "you have reached the maximum number of requests or actions allowed within a certain time frame"
)

type MiddlewareRTL struct {
	RTL *rtl.RTL
}

func NewMiddlewareRTL(rtl *rtl.RTL) *MiddlewareRTL {
	return &MiddlewareRTL{
		RTL: rtl,
	}
}

func (m *MiddlewareRTL) RateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parameter := r.Header.Get("API_KEY")
		if parameter == "" {
			parameter = r.RemoteAddr
		}
		validate, err := m.RTL.Validate(r.Context(), parameter)
		if err != nil || !validate {
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(errorRT))
			return
		}
		next.ServeHTTP(w, r)
	})
}
