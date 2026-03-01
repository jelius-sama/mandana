package api

import (
    "mandana/api/handler"
    "net/http"
)

func Router() *http.ServeMux {
    var mux *http.ServeMux = http.NewServeMux()

    mux.HandleFunc("GET /api/stats/{$}", handler.HandleStats)

    return mux
}

