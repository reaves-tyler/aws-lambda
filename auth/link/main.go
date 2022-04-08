package main

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt"
)

type Claims struct {
	// foo string
	// nbf int64
}

type Token struct {
	Raw string // The raw token.  Populated when you Parse a token
	// Method    SigningMethod          // The signing method used or to be used
	Header map[string]interface{} // The first segment of the token
	// Claims    Claims                 // The second segment of the token
	Signature string // The third segment of the token.  Populated when you Parse a token
	Valid     bool   // Is the token valid?  Populated when you Parse/Verify a token
}

type AuthBody struct {
	Url       string            `json:"url"`
	UrlParams map[string]string `json:"urlParams"`
}

type AuthResponseBody struct {
	Link string `json:"link"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var authBody AuthBody
	json.Unmarshal([]byte(request.Body), &authBody)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss": "Trader Interactive",
		"sub": "Auth Token",
		"aud": "some-realm",
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"nbf": time.Now().Unix(),
		"iat": time.Now().Unix(),
		"url": authBody.Url,
	})

	// Sign and get the complete encoded token as a string using the secret
	rawDecodedText, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRET"))

	if err != nil {
		panic(err)
	}

	tokenString, err := token.SignedString(rawDecodedText)

	if err != nil {
		log.Fatal(err)
	}

	tokenString = "token=" + tokenString

	var responseLink AuthResponseBody = AuthResponseBody{
		Link: authBody.Url + "?" + tokenString,
	}

	res, err := json.Marshal(responseLink)

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
