package handler

import (
    "encoding/json"
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

