package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"google.golang.org/api/option"
)

// Firestore クライアントを作成
func createFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	sa := option.WithCredentialsFile("serviceAccountKey.json")
	client, err := firestore.NewClient(ctx, "livegift-37bc2", sa)
	if err != nil {
		return nil, fmt.Errorf("Firestore クライアントの作成に失敗: %v", err)
	}
	return client, nil
}

// Firestore からデータを取得する Lambda ハンドラー
func getMessagesFromFirestore(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	client, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("Firestore クライアントエラー: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}
	defer client.Close()

	// Firestore の "messages" コレクションからデータを取得
	iter := client.Collection("messages").Documents(ctx)
	var messages []map[string]interface{}

	for {
		doc, err := iter.Next()
		if err != nil {
			break
		}
		messages = append(messages, doc.Data())
	}

	// JSON に変換
	jsonData, err := json.Marshal(messages)
	if err != nil {
		log.Printf("JSON 変換エラー: %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Failed to parse data"}, nil
	}

	return events.APIGatewayV2HTTPResponse{
		StatusCode: 200,
		Body:       string(jsonData),
		Headers:    map[string]string{"Content-Type": "application/json"},
	}, nil
}

func addMessageToFirestore(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	client, err := createFirestoreClient(ctx)
	if err != nil {
		log.Printf("Firestore クライアントエラー：%v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Internal Server Error"}, nil
	}
	defer client.Close()

	// リクエストのJSONをパース
	var requestBody struct {
		Message string `json: "message"`
	}
	err = json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		log.Printf("JSON パースエラー： %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 400, Body: "Invalid JSON"}, nil
	}

	// Firestoreにデータを追加
	_, _, err = client.Collection("messages").Add(ctx, map[string]interface{}{
		"message":   requestBody.Message,
		"createdAt": time.Now(),
	})
	if err != nil {
		log.Printf("Firestore 書き込みエラー： %v", err)
		return events.APIGatewayV2HTTPResponse{StatusCode: 500, Body: "Failed to save message"}, nil
	}
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 201,
		Body:       "Message added successfully",
	}, nil
}

func main() {
	lambda.Start(getMessagesFromFirestore)
}
