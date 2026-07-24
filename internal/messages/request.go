package messages

import (
	"encoding/base64"
	"io"
	"net/http"
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
	Body                            any                 `json:"body"`
	IsBase64Encoded                 bool                `json:"isBase64Encoded"`
}

func ParseRequest(r *http.Request) (Request, error) {
	body, isBase64Encoded, err := readBody(r)
	if err != nil {
		return Request{}, err
	}
	headers, multiValueHeaders := flatten(r.Header)
	query, multiValueQuery := flatten(r.URL.Query())

	return Request{
		Resource:                        r.URL.Path,
		Path:                            r.URL.Path,
		HTTPMethod:                      r.Method,
		Headers:                         headers,
		MultiValueHeaders:               multiValueHeaders,
		QueryStringParameters:           query,
		MultiValueQueryStringParameters: multiValueQuery,
		Body:                            body,
		IsBase64Encoded:                 isBase64Encoded,
	}, nil
}

func flatten(values map[string][]string) (map[string]string, map[string][]string) {
	single := make(map[string]string, len(values))
	for k, v := range values {
		if len(v) > 0 {
			single[k] = v[len(v)-1]
		}
	}
	return single, values
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
