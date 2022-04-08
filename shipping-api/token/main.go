package main

import (
	"encoding/json"
	"log"
	"shipping-api/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	data, err := call()

	if err != nil {
		log.Fatal(err)
	}

	d, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(d),
		StatusCode: 200,
	}, nil
}

type User struct {
	Username    string   `json:"username"`
	Roles       []string `json:"roles"`
	TokenType   string   `json:"token_type"`
	AccessToken string   `json:"access_token"`
}

func call() (User, error) {
	body, err := utils.MakeRequest("POST", "login")

	if err != nil {
		log.Fatal(err)
		return User{}, err
	}

	var user User
	json.Unmarshal([]byte(string(body)), &user)

	return user, nil
}

func main() {
	lambda.Start(Handler)
}
