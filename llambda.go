package llambda

import (
	"fmt"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/dbrun3/llambda/internal"
)

// Start takes a handler and talks to an internal Lambda endpoint to pass requests to the handler. If the
// handler does not match one of the supported types an appropriate error message will be returned to the caller.
// Start blocks, and does not return after being called.
//
// When run in a non-Lambda environment, a local http server will start at LAMBDA_SERVER_PORT (default 8080)
// to proxy requests to the internal Lambda which will always be available on the following port (default 8081).
func Start(handler any) {
	StartWithOptions(handler)
}

// StartWithOptions is the same as Start after the application of any handler options specified
func StartWithOptions(handler any, options ...lambda.Option) {
	if os.Getenv("AWS_LAMBDA_FUNCTION_NAME") == "" {
		port := 8080
		if setPort := os.Getenv("LAMBDA_SERVER_PORT"); setPort != "" {
			port, _ = strconv.Atoi(setPort)
		}
		os.Setenv("_LAMBDA_SERVER_PORT", fmt.Sprintf("%d", port+1))
		http.Handle("/", &internal.Proxy{
			LambdaClient: sync.OnceValue(func() *rpc.Client {
				c, _ := rpc.Dial("tcp", fmt.Sprintf(":%d", port+1))
				return c
			}),
		})
		go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}
	lambda.StartWithOptions(handler, options...)
}
