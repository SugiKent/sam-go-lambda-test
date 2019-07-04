.PHONY: deps clean build

deps:
	go get -u ./...

clean:
	rm -rf ./hello-world/hello-world

build:
	# GOOS=linux GOARCH=amd64 go build -o hello-world/hello-world ./hello-world
	# GOOS=linux GOARCH=amd64 go build -o score-register/score-register ./score-register
	GOOS=linux GOARCH=amd64 go build -o timelog-register/timelog-register ./timelog-register
	GOOS=linux GOARCH=amd64 go build -o timelog-fetcher/timelog-fetcher ./timelog-fetcher

package:
	sam package --template-file template.yaml --output-template-file output-template.yaml --s3-bucket go-sam-template-store --profile sugita

deploy:
	sam deploy --template-file output-template.yaml --stack-name go-sam-template-store --capabilities CAPABILITY_IAM --profile sugita

dynamodb:
	aws dynamodb create-table --table-name Score --attribute-definitions AttributeName=PersonID,AttributeType=S AttributeName=TestID,AttributeType=S --key-schema AttributeName=PersonID,KeyType=HASH AttributeName=TestID,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --profile sugita

dynamo-timelog:
	aws dynamodb create-table --table-name Timelog --attribute-definitions AttributeName=TimeID,AttributeType=S AttributeName=StartAt,AttributeType=S --key-schema AttributeName=TimeID,KeyType=HASH AttributeName=StartAt,KeyType=RANGE --provisioned-throughput ReadCapacityUnits=1,WriteCapacityUnits=1 --profile sugita
