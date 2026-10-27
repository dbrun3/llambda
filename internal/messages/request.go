package messages

import (
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"
)

type Request struct {
	Resource                        string              `json:"resource"`
	Path                            string              `json:"path"`
	HTTPMethod                      string              `json:"httpMethod"`
	Headers                         map[string]string   `json:"headers"`
	MultiValueHeaders               map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters           map[string]string   `json:"queryStringParameters"`
	MultiValueQueryStringParameters map[string][]string `json:"multiValueQueryStringParameters"`
	Body                            string              `json:"body"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded"`
}

func ParseRequest(r *http.Request) (Request, error) {
	body, isBase64Encoded, err := readBody(r)
	if err != nil {
		return Request{}, err
	}

	headers := lower(r.Header)
	queryParams := r.URL.Query()
	if len(queryParams) == 0 {
		queryParams = nil
	}

	return Request{
		Resource:                        r.URL.Path,
		Path:                            r.URL.Path,
		HTTPMethod:                      r.Method,
		Headers:                         flatten(headers),
		MultiValueHeaders:               headers,
		QueryStringParameters:           flatten(queryParams),
		MultiValueQueryStringParameters: queryParams,
		Body:                            body,
		IsBase64Encoded:                 isBase64Encoded,
	}, nil
}

func lower(values map[string][]string) map[string][]string {
	lowered := make(map[string][]string, len(values))
	for k, v := range values {
		if len(v) > 0 {
			lowered[strings.ToLower(k)] = v
		}
	}
	return lowered
}

func flatten(values map[string][]string) map[string]string {
	if len(values) == 0 {
		return nil
	}
	single := make(map[string]string, len(values))
	for k, v := range values {
		if len(v) > 0 {
			single[k] = v[len(v)-1]
		}
	}
	return single
}

func readBody(r *http.Request) (string, bool, error) {
	if r.Body == nil {
		return "", false, nil
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", false, err
	}
	if !utf8.Valid(body) {
		return base64.StdEncoding.EncodeToString(body), true, nil
	}
	return string(body), false, nil
}
