package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-sdk-go/aws"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	reqBody := request.Body
	fmt.Println(reqBody)
	jsonBytes := ([]byte)(reqBody)
	timelogReq := new(TimelogRequest)
	if err := json.Unmarshal(jsonBytes, timelogReq); err != nil {
		fmt.Println("[ERROR]", err)
	}

	timeID := timelogReq.TimeID
	startAt := timelogReq.StartAt
	endAt := timelogReq.EndAt

	timelog := Timelog{
		TimeID:  timeID,
		StartAt: startAt,
		EndAt:   endAt,
	}

	av, err := dynamodbattribute.MarshalMap(timelog)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}

	session, err := session.NewSession()
	conn := dynamodb.New(session)
	param, err := conn.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("Timelog"),
		Item:      av,
	})
	if err != nil {
		fmt.Println("[ERROR]", err)
	}
	fmt.Println(param)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		Headers: map[string]string{
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Headers": "Content-Type",
			"Content-Type":                 "application/json",
		},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type TimelogRequest struct {
	TimeID  string `json:"TimeID"`
	StartAt string `json:"StartAt"`
	EndAt   string `json:"EndAt"`
}

type Timelog struct {
	TimeID  string
	StartAt string
	EndAt   string
}
