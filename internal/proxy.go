package internal

import (
	"encoding/json"
	"net/http"
	"net/rpc"
	"time"

	lambdaMessages "github.com/aws/aws-lambda-go/lambda/messages"
	"github.com/dbrun3/llambda/internal/messages"
)

const functioninvokeRPC = "Function.Invoke"

type Proxy struct {
	LambdaClient func() *rpc.Client
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	invokeRequest, err := createInvokeRequest(r)
	if err != nil {
		messages.GatewayError(err.Error(), w)
		return
	}

	var invokeResponse lambdaMessages.InvokeResponse
	p.LambdaClient().Call(functioninvokeRPC, invokeRequest, &invokeResponse)
	if invokeResponse.Error != nil {
		messages.GatewayError(invokeResponse.Error.Message, w)
		return
	}

	var response messages.Response
	err = json.Unmarshal(invokeResponse.Payload, &response)
	if err != nil {
		messages.GatewayError(err.Error(), w)
		return
	}
	response.Write(w)
}

// https://github.com/djhworld/go-lambda-invoke/blob/master/golambdainvoke/invoke.go
func createInvokeRequest(r *http.Request) (*lambdaMessages.InvokeRequest, error) {
	request, err := messages.ParseRequest(r)
	if err != nil {
		return nil, err
	}

	payloadEncoded, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}

	t := time.Now()
	deadline := &lambdaMessages.InvokeRequest_Timestamp{
		Seconds: int64(t.Unix()),
		Nanos:   int64(t.Nanosecond()),
	}

	return &lambdaMessages.InvokeRequest{
		Payload:   payloadEncoded,
		RequestId: "0",
		Deadline:  *deadline,
	}, nil
}
