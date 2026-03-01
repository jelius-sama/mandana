package middleware

import (
    "net/http"

    "github.com/jelius-sama/logger"
    "runtime/debug"
)

// TODO: Fix the issue where the response might already be half-written before the server panicked.
func RecoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                logger.TimedError("Encountered a panic, returning 500 to client and recovering the server!")
                logger.TimedError(
                    "panic recovered",
                    "\n\terror", err,
                    "\n\tstack", string(debug.Stack()),
                )
                w.Header().Set("Content-Type", "application/json")
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte(`{"code":"500", "error": "Internal Server Error"}`))
            }
        }()
        next.ServeHTTP(w, r)
    })
}

