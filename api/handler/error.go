package handler

import (
    "encoding/json"
    "fmt"
    "mandana/types"
    "net/http"

    "github.com/jelius-sama/logger"
)

func HttpError(w http.ResponseWriter, err error, code int) {
    w.Header().Set("Content-Type", "application/json")

    if err = json.NewEncoder(w).Encode(types.HTTPError{
        Code:    code,
        Error:   http.StatusText(code),
        Message: err.Error(),
    }); err != nil {
        logger.Panic("Failed to encode http error:", err.Error())
    }
}

func HTTPPlaceholder(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusNotImplemented)

    if _, err := fmt.Fprintf(w, `{"code":"501", "error": "Not Implemented", "message": "TODO: Implement %c%s%c route."}`, '`', r.URL.Path, '`'); err != nil {
        logger.Panic("Failed to write http response:", err.Error())
    }
}

