package messages

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Payload struct {
	UserId  int    `json:"userId"`
	Message string `json:"message"`
}

func TestParseRequest(t *testing.T) {
	assert := assert.New(t)

	t.Run("headers are lowercase", func(t *testing.T) {
		headerReq := http.Request{
			Header: map[string][]string{"Authorization": {"Bearer key"}},
			URL:    &url.URL{},
		}

		headerExpect := Request{
			Headers:           map[string]string{"authorization": "Bearer key"},
			MultiValueHeaders: map[string][]string{"authorization": {"Bearer key"}},
		}
		headerResult, err := ParseRequest(&headerReq)
		assert.NoError(err)
		assert.Equal(headerExpect, headerResult)
	})

	t.Run("multival headers,params parsed and collapsed", func(t *testing.T) {
		multiValReq := http.Request{
			URL: &url.URL{
				RawQuery: "world=true&world=false&name=user",
			},
			Header: map[string][]string{"authorization": {"Bearer key"}, "app": {"hello-world", "hello-friend"}},
		}

		multiValExpect := Request{
			QueryStringParameters:           map[string]string{"world": "false", "name": "user"},
			MultiValueQueryStringParameters: map[string][]string{"world": {"true", "false"}, "name": {"user"}},
			Headers:                         map[string]string{"authorization": "Bearer key", "app": "hello-friend"},
			MultiValueHeaders:               map[string][]string{"authorization": {"Bearer key"}, "app": {"hello-world", "hello-friend"}},
		}
		multiValResult, err := ParseRequest(&multiValReq)
		assert.NoError(err)
		assert.Equal(multiValExpect, multiValResult)
	})

	t.Run("can parse json body", func(t *testing.T) {
		expectedPayload := Payload{UserId: 123, Message: "Hello"}
		payload, _ := json.Marshal(expectedPayload)
		bodyReq := http.Request{
			Body: io.NopCloser(bytes.NewReader(payload)),
			URL:  &url.URL{},
		}
		bodyRes, err := ParseRequest(&bodyReq)
		assert.NoError(err)

		var receivedPayload Payload
		json.Unmarshal([]byte(bodyRes.Body), &receivedPayload)
		assert.Equal(expectedPayload, receivedPayload)
	})

	t.Run("can parse b64 body", func(t *testing.T) {
		expectedPayload := Payload{UserId: 123, Message: "Hello"}
		payload, _ := json.Marshal(expectedPayload)

		buf := new(bytes.Buffer)
		binary.Write(buf, binary.BigEndian, uint32(0xFFFFFFFF))
		binary.Write(buf, binary.BigEndian, payload)

		bodyReq := http.Request{
			Body: io.NopCloser(bytes.NewReader(buf.Bytes())),
			URL:  &url.URL{},
		}
		bodyRes, err := ParseRequest(&bodyReq)
		assert.NoError(err)

		b64Payload, err := base64.StdEncoding.DecodeString(bodyRes.Body)
		assert.NoError(err)

		var receivedPayload Payload
		json.Unmarshal(b64Payload[4:], &receivedPayload)

		assert.Equal(expectedPayload, receivedPayload)
	})

}
