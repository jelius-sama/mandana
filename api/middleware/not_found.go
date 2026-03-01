package middleware

import "net/http"

func NotFound(mux *http.ServeMux) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var _, pattern = mux.Handler(r)

        if pattern == "" {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusNotFound)
            w.Write([]byte(`{"status":"404", "message": "Not Found"}`))
            return
        }

        mux.ServeHTTP(w, r)
    })
}

