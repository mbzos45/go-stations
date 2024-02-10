package middleware

import (
	"context"
	"fmt"
	"github.com/mileusna/useragent"
	"net/http"
)

const UAKey = "user_agent"

func SetUA(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ua := useragent.Parse(r.UserAgent())
	ctx := context.WithValue(r.Context(), UAKey, ua)
	next(rw, r.WithContext(ctx))
}

func GetClientOS(ctx context.Context) (string, error) {
	ua, ok := ctx.Value(UAKey).(useragent.UserAgent)
	if !ok {
		return "", fmt.Errorf("userAgent not found")
	}
	return ua.OS, nil
}
