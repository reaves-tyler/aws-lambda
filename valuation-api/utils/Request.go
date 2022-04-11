package utils

import (
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

func MakeRequest(method string, path string) ([]byte, int, error) {
	req, err := http.NewRequest(method, os.Getenv("BASEURL")+path, nil)

	if err != nil {
		log.Fatal(err)
		return nil, http.StatusInternalServerError, err
	}

	AppendJDPowerAuth(req)

	body, statusCode, err := DoRequest(req)

	if err != nil {
		return nil, statusCode, err
	}

	return body, statusCode, nil
}

func AppendJDPowerAuth(req *http.Request) *http.Request {
	req.Header.Add("UserName", os.Getenv("USERNAME"))
	req.Header.Add("Password", os.Getenv("PASSWORD"))

	return req
}

func HandleError(err error, statusCode int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		Body:       err.Error(),
		StatusCode: statusCode,
	}, nil
}
