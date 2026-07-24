package messages

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
)

type Response struct {
	StatusCode        int                 `json:"statusCode"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Body              any                 `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}

func (r Response) Write(w http.ResponseWriter) {
	body, err := r.body()
	if err != nil {
		GatewayError(err.Error(), w)
		return
	}
	for k, v := range r.Headers {
		w.Header().Set(k, v)
	}
	for k, vs := range r.MultiValueHeaders {
		for _, v := range vs {
			w.Header().Add(k, v)
		}
	}
	w.WriteHeader(r.StatusCode)
	w.Write(body)
}

func (r Response) body() ([]byte, error) {
	if !r.IsBase64Encoded {
		return json.Marshal(r.Body)
	}
	s, ok := r.Body.(string)
	if !ok {
		return nil, fmt.Errorf("isBase64Encoded set but body is %T, want string", r.Body)
	}
	return base64.StdEncoding.Strict().DecodeString(s)
}

func GatewayError(msg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusBadGateway)
	fmt.Fprintf(w, "<html><body><h1>502 Bad Gateway</h1><p>%s</p></body></html>", html.EscapeString(msg))
}
