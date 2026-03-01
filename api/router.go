package api

import (
    "fmt"
    "mandana/api/handler"
    "net/http"
    "strings"

    "github.com/jelius-sama/logger"
)

type HTTPMethod int

const (
    MethodGET HTTPMethod = iota
    MethodPOST
    MethodPATCH
    MethodPUT
    MethodDELETE
)

func (hm HTTPMethod) String() string {
    switch hm {
    case MethodGET:
        return "GET"
    case MethodPOST:
        return "POST"
    case MethodPATCH:
        return "PATCH"
    case MethodPUT:
        return "PUT"
    case MethodDELETE:
        return "DELETE"
    }

    // TODO: Handle other HTTP Methods
    logger.Panic("Unreachable")
    return "GET"
}

// Generates and returns an absolute path
func absPath(path string, method HTTPMethod) string {
    return fmt.Sprintf(
        "%s /api/%s/{$}",
        method,
        func() string {
            var cleaned, _ = strings.CutPrefix(path, "/")
            return cleaned
        }(),
    )
}

// Generates and returns a generic path
func genPath(path string, method HTTPMethod) string {
    return fmt.Sprintf(
        "%s /api/%s/",
        method,
        func() string {
            var cleaned, _ = strings.CutPrefix(path, "/")
            return cleaned
        }(),
    )
}

func Router() *http.ServeMux {
    var mux *http.ServeMux = http.NewServeMux()

    mux.HandleFunc(absPath("/get/all", MethodGET), handler.HTTPPlaceholder)
    mux.HandleFunc(genPath("get", MethodGET), handler.HTTPPlaceholder)

    mux.HandleFunc(absPath("stats", MethodGET), handler.HandleStats)

    return mux
}

