package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
)

// ハンドラー関数
func handler(ctx context.Context) (string, error) {
	return "Hello from AWS Lambda!", nil
}

func main() {
	fmt.Println("AWS Lambda is running...")
	lambda.Start(handler)
}
