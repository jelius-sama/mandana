package handler

import (
    "net/http"

    "github.com/jelius-sama/logger"
)

func HandleGetPanel(w http.ResponseWriter, r *http.Request) {
    // TODO: implement panel loading
    // Expected route: /api/panel/%s?page=%d

    const placeholder = "/assets/resource/placeholder.jpg"

    logger.Debug("TODO: HandleGetPanel not implemented, redirecting to placeholder asset")

    w.Header().Set("Location", placeholder)
    w.WriteHeader(http.StatusTemporaryRedirect)
}

