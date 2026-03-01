package middleware

import (
    "net/http"

    "github.com/jelius-sama/logger"
)

func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        logger.TimedInfo(r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

