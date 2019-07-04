package main

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	endAt := "false"
	session, err := session.NewSession()
	conn := dynamodb.New(session)

	param, err := conn.Query(&dynamodb.QueryInput{
		TableName: aws.String("Timelog"),
		ExpressionAttributeNames: map[string]*string{
			"#EndAt": aws.String("EndAt"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":endAt": {
				S: aws.String(endAt),
			},
		},
		KeyConditionExpression: aws.String("#EndAt = :endAt"),
		IndexName:              aws.String("EndAt-StartAt-index"),
	})

	if err != nil {
		fmt.Println("[ERROR]", err)
	}
	fmt.Println(param)

	timelogs := make([]*TimeLogRes, 0)
	if err := dynamodbattribute.UnmarshalListOfMaps(param.Items, &timelogs); err != nil {
		fmt.Println("[ERROR]", err)
	}
	jsonBytes, _ := json.Marshal(timelogs)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		Headers: map[string]string{
			"Access-Control-Allow-Origin": "*",
		},
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}

type TimeLogRes struct {
	TimeID  string `json:"timeID"`
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
}
