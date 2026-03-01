package types

type HTTPError struct {
    Code    int    `json:"code"`
    Error   string `json:"error"`
    Message string `json:"message"`
}

