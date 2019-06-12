package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(func(ctx context.Context, in interface{}) error {
		fmt.Printf("%+v\n", in)
		return nil
	})
}
