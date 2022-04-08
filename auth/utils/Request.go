package utils

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
)

func DoRequest(req *http.Request) ([]byte, int, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return nil, res.StatusCode, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
		return nil, http.StatusInternalServerError, err
	}

	return body, res.StatusCode, nil
}

func MakeTokenRequest() ([]byte, int, error) {
	req, err := http.NewRequest("POST", os.Getenv("NEBULOUS_URL")+"token", nil)

	if err != nil {
		log.Fatal(err)
		return nil, http.StatusInternalServerError, err
	}

	body, statusCode, err := DoRequest(req)

	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func MakeRequest(method string, endpoint string, body []byte, token string) ([]byte, int, error) {
	req, err := http.NewRequest(method, endpoint, bytes.NewBuffer(body))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	if err != nil {
		log.Fatal(err)
		return nil, http.StatusInternalServerError, err
	}

	body, statusCode, err := DoRequest(req)

	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func HandleError(err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: statusCode,
	}, nil
}
