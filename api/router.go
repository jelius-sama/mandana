package api

import (
    "fmt"
    "mandana/api/handler"
    "mandana/client/assets"
    "mandana/client/pages"
    "mandana/types"
    "mime"
    "net/http"
    "strconv"
    "strings"

    "github.com/google/uuid"
    "github.com/jelius-sama/logger"
)

func init() {
    // Force correct MIME types for web assets (handles some edge cases)
    var jsMimeErr error = mime.AddExtensionType(".js", "application/javascript")
    if jsMimeErr != nil {
        logger.Error("Failed to set mime for js file types!")
    }
    var cssMimeErr error = mime.AddExtensionType(".css", "text/css")
    if cssMimeErr != nil {
        logger.Error("Failed to set mime for css file types!")
    }
}

type RouteT uint8
type HTTPMethod uint8

const (
    MethodGET HTTPMethod = iota
    MethodPOST
    MethodPATCH
    MethodPUT
    MethodDELETE
)

const (
    RouteAPI RouteT = iota
    RoutePage
    RouteAsset
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

func (rt RouteT) String() string {
    switch rt {
    case RouteAPI:
        return "/api"
    case RouteAsset:
        return "/assets"
    case RoutePage:
        return ""
    }
    logger.Panic("Unreachable")
    return ""
}

// Generates and returns an absolute path
func absPath(path string, method HTTPMethod, routeType RouteT) string {
    path, _ = strings.CutPrefix(path, "/")
    path, _ = strings.CutSuffix(path, "/")
    if path != "" {
        path = fmt.Sprintf("/%s", path)
    }
    return fmt.Sprintf("%s %s%s/{$}", method, routeType, path)
}

// Generates and returns a generic path
func genPath(path string, method HTTPMethod, routeType RouteT, params ...string) string {
    path, _ = strings.CutPrefix(path, "/")
    path, _ = strings.CutSuffix(path, "/")
    if path != "" {
        path = fmt.Sprintf("/%s", path)
    }
    if len(params) > 0 {
        return fmt.Sprintf("%s %s%s/%s", method, routeType, path, strings.Join(params, ""))
    } else {
        return fmt.Sprintf("%s %s%s/", method, routeType, path)
    }
}

func Router() *http.ServeMux {
    var mux *http.ServeMux = http.NewServeMux()

    mux.HandleFunc(absPath("/", MethodGET, RoutePage), func(w http.ResponseWriter, r *http.Request) {
        handler.Page(handler.PageT{
            W:    w,
            R:    r,
            Page: pages.Home(),
        })
    })

    mux.HandleFunc(genPath("yomu", MethodGET, RoutePage, "{id}"), func(w http.ResponseWriter, r *http.Request) {
        var uuid, parseErr = uuid.Parse(r.PathValue("id"))

        if parseErr != nil {
            handler.Page(handler.PageT{
                W:    w,
                R:    r,
                Page: pages.NotFound(),
            })
            return
        }

        var sakuhin *types.Sakuhin = &types.Sakuhin{
            ID: uuid,
            Title: types.StringWithLang{
                English: logger.StringPtr("Placeholder Title"),
            },
            CoverArts: []string{"/assets/resource/placeholder.jpg"},
            PageCount: 10,
        }

        pageStr := r.URL.Query().Get("page")
        if len(pageStr) == 0 {
            handler.Page(handler.PageT{
                W:    w,
                R:    r,
                Page: pages.Sakuhin(sakuhin),
            })
            return
        }

        page, convertErr := strconv.Atoi(pageStr)
        if page < 1 || convertErr != nil {
            page = 1
        }

        handler.Page(handler.PageT{
            W:    w,
            R:    r,
            Page: pages.Yomu(sakuhin, uint16(page)),
        })
    })

    mux.HandleFunc(genPath("/", MethodGET, RoutePage), func(w http.ResponseWriter, r *http.Request) {
        handler.Page(handler.PageT{
            W:    w,
            R:    r,
            Page: pages.NotFound(),
        })
    })

    mux.HandleFunc(absPath("/sakuhin/all", MethodGET, RouteAPI), handler.HTTPPlaceholder)
    mux.HandleFunc(genPath("sakuhin", MethodGET, RouteAPI), handler.HTTPPlaceholder)

    mux.HandleFunc(genPath("panel", MethodGET, RouteAPI, "{id}"), handler.HandleGetPanel)

    mux.HandleFunc(absPath("stats", MethodGET, RouteAPI), handler.HandleStats)

    mux.Handle(genPath("css", MethodGET, RouteAsset), http.StripPrefix("/assets/",
        http.FileServer(http.FS(assets.Assets))))

    mux.Handle(genPath("js", MethodGET, RouteAsset), http.StripPrefix("/assets/",
        http.FileServer(http.FS(assets.Assets))))

    mux.Handle(genPath("fonts", MethodGET, RouteAsset), http.StripPrefix("/assets/",
        http.FileServer(http.FS(assets.Assets))))

    mux.Handle(genPath("resource", MethodGET, RouteAsset), http.StripPrefix("/assets/",
        http.FileServer(http.FS(assets.Assets))))

    // utils.SetupScriptRoutes(mux, os.Getenv("IS_PROD") == "FALSE")

    return mux
}

