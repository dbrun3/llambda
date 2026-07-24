package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/dbrun3/llambda"
)

// handler returns "hello <path>!"
func handler(ctx context.Context, r events.ALBTargetGroupRequest) (events.ALBTargetGroupResponse, error) {
	return events.ALBTargetGroupResponse{
		StatusCode: http.StatusOK,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		Body:       fmt.Sprintf("hello %s!", r.Path[1:]),
	}, nil
}

func main() {
	llambda.Start(handler)
}
