package main

import (
	"context"
	"fmt"
	"log"

	"cloud.google.com/go/firestore"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
)

// Firestore クライアントを作成
func createFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	client, err := firestore.NewClient(ctx, "YOUR_FIREBASE_PROJECT_ID", sa)
	if err != nil {
		return nil, fmt.Errorf("Firestore クライアントの作成に失敗: %v", err)
	}
	return client, nil
}

// Firestore にデータを追加する Lambda ハンドラー
func addDataToFirestore(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	client, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("Firestore クライアントエラー: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}
	defer client.Close()

	_, _, err = client.Collection("messages").Add(ctx, map[string]interface{}{
		"message": "Hello from AWS Lambda & Firestore!",
	})
	if err != nil {
		log.Printf("Firestore 書き込みエラー: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Failed to add data"}, nil
	}

	return events.APIGatewayV2HTTPResponse{StatusCode: 200, Body: "Data added to Firestore"}, nil
}

func main() {
	lambda.Start(addDataToFirestore)
}
