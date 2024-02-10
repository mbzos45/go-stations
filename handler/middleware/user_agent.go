package middleware

import (
	"context"
	"fmt"
	"github.com/mileusna/useragent"
	"net/http"
)

const clientOSKey = "client_OS"

func SetClientOS(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ua := useragent.Parse(r.UserAgent())
	ctx := context.WithValue(r.Context(), clientOSKey, ua.OS)
	next(rw, r.WithContext(ctx))
}

func GetClientOS(ctx context.Context) (string, error) {
	os, ok := ctx.Value(clientOSKey).(string)
	if !ok {
		return "", fmt.Errorf("clientOS not found")
	}
	return os, nil
}
