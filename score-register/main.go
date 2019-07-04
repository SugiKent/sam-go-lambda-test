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
	personReq := new(PersonRequest)
	if err := json.Unmarshal(jsonBytes, personReq); err != nil {
		fmt.Println("[ERROR]", err)
	}

	personID := personReq.PersonID
	personName := personReq.PersonName
	testID := personReq.TestID
	score := personReq.Score
	passingMark := false
	if score >= 80 {
		passingMark = true
	}

	// DynamoDBへ永続化
	person := Person{
		PersonID:    personID,
		PersonName:  personName,
		TestID:      testID,
		Score:       score,
		PassingMark: passingMark,
	}
	av, err := dynamodbattribute.MarshalMap(person)
	if err != nil {
		fmt.Println("[ERROR]", err)
	}

	session, err := session.NewSession()
	conn := dynamodb.New(session)
	param, err := conn.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String("Score"),
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

type PersonRequest struct {
	PersonID   string `json:"personID"`
	PersonName string `json:"personName"`
	TestID     string `json:"testID"`
	Score      int    `json:"score"`
}

type Person struct {
	PersonID    string
	PersonName  string
	TestID      string
	Score       int
	PassingMark bool
}
