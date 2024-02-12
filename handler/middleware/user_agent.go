package middleware

import (
	"context"
	"fmt"
	"github.com/mileusna/useragent"
	"net/http"
)

const UAKey = "user_agent"

func SetUA(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ua := useragent.Parse(r.UserAgent())
		ctx := context.WithValue(r.Context(), UAKey, ua)
		h.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func GetClientOS(ctx context.Context) (string, error) {
	ua, ok := ctx.Value(UAKey).(useragent.UserAgent)
	if !ok {
		return "", fmt.Errorf("userAgent not found")
	}
	return ua.OS, nil
}
