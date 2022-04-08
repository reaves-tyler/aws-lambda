package utils

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

func DoRequest(req *http.Request) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return body, nil
}

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func MakeRequest(method string, path string) ([]byte, error) {
	authData := Auth{
		Username: os.Getenv("USERNAME"),
		Password: os.Getenv("PASSWORD"),
	}

	authJson, err := json.Marshal(authData)

	if err != nil {
		log.Fatal(err)
	}

	req, err := http.NewRequest(method, os.Getenv("BASEURL")+path, bytes.NewBuffer(authJson))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	body, err := DoRequest(req)

	if err != nil {
		return nil, err
	}

	return body, nil
}
