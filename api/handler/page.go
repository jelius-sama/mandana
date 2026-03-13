package handler

import (
    "mandana/client"
    "net/http"

    "github.com/a-h/templ"
)

type PageT struct {
    W    http.ResponseWriter
    R    *http.Request
    Page templ.Component
}

func Page(props PageT) {
    if props.R.Header.Get("HX-Request") == "true" {
        err := props.Page.Render(props.R.Context(), props.W)
        if err != nil {
            panic(err)
        }
    } else {
        ctx := templ.WithChildren(props.R.Context(), props.Page)
        err := client.Layout().Render(ctx, props.W)
        if err != nil {
            panic(err)
        }
    }
}

