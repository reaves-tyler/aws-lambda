package main

import (
	"auth/utils"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/golang-jwt/jwt"
)

type Token struct {
	Raw string // The raw token.  Populated when you Parse a token
	// Method    SigningMethod          // The signing method used or to be used
	Header map[string]interface{} // The first segment of the token
	// Claims    Claims                 // The second segment of the token
	Signature string // The third segment of the token.  Populated when you Parse a token
	Valid     bool   // Is the token valid?  Populated when you Parse/Verify a token
}

type AuthBody struct {
	Email     string            `json:"email"`
	Url       string            `json:"url"`
	UrlParams map[string]string `json:"urlParams"`
	Realm     string            `json:"realm"`
}

type AuthResponseBody struct {
	Link string `json:"link"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var authBody AuthBody
	json.Unmarshal([]byte(request.Body), &authBody)

	err := validateRequestBody(&authBody)

	if err != nil {
		return utils.HandleError(err, http.StatusBadRequest)
	}

	var CustomerIDForToken string

	customerID, statusCode, err := utils.CheckCustomerExists(authBody.Email)
	CustomerIDForToken = customerID

	if err != nil {
		log.Fatal(err)
	}

	if statusCode >= 400 {
		return utils.HandleError(err, statusCode)
	}

	if CustomerIDForToken == "" {
		customer, statusCode, err := utils.CreateCustomer(authBody.Email)

		if err != nil {
			log.Fatal(err)
		}

		if statusCode >= 400 {
			return utils.HandleError(err, statusCode)
		}

		CustomerIDForToken = customer.Result.ID
	}

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"iss":   "Trader Interactive",
		"sub":   CustomerIDForToken,
		"aud":   authBody.Realm,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"nbf":   time.Now().Unix(),
		"iat":   time.Now().Unix(),
		"url":   authBody.Url,
		"email": authBody.Email,
	})

	rawDecodedText, err := base64.StdEncoding.DecodeString(os.Getenv("JWT_SECRET"))

	if err != nil {
		log.Fatal(err)
	}

	// https://pkg.go.dev/encoding/pem#example-Decode
	block, rest := pem.Decode(rawDecodedText)
	if block == nil {
		log.Fatal("failed to decode PEM block containing private key", rest)
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)

	if err != nil {
		log.Print("failed to parse RSA public Key from PEM")
		log.Fatal(err)
	}

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(priv)

	if err != nil {
		log.Fatal(err)
	}

	url := addQueryParams(&authBody, &tokenString)

	var responseLink AuthResponseBody = AuthResponseBody{
		Link: url.String(),
	}

	res, err := json.Marshal(responseLink)

	if err != nil {
		log.Fatal(err)
	}

	utils.SendEmail(authBody.Email, url.String(), authBody.Realm)

	return events.APIGatewayProxyResponse{
		Body:       string(res),
		StatusCode: 200,
	}, nil
}

func addQueryParams(body *AuthBody, token *string) *url.URL {
	url, err := url.Parse(body.Url)
	if err != nil {
		log.Fatal(err)
	}

	queryParams := url.Query()
	queryParams.Add("token", *token)

	for key, value := range (*body).UrlParams {
		queryParams.Add(key, value)
	}

	url.RawQuery = queryParams.Encode()

	return url
}

func validateRequestBody(body *AuthBody) error {
	if body.Email == "" {
		return errors.New("email is required")
	}

	if body.Url == "" {
		return errors.New("url is required")
	}

	if body.Realm == "" {
		return errors.New("realm is required")
	}

	return nil
}

func main() {
	lambda.Start(Handler)
}
