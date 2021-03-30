package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Lolik232/luximapp-api/internal/context_keys"
)

func DeviceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			userAgent := r.UserAgent()
			fmt.Println(userAgent)
			// re := regexp.MustCompile(`\(([a-zA-Z0-9;\. _+-]+)\)`)
			device := userAgent
			ctx := context.WithValue(r.Context(), context_keys.DeviceInfoContextKey, device)
			fmt.Printf("Device connected:  %s\n", device)
			next.ServeHTTP(w, r.WithContext(ctx))
		},
	)
}
