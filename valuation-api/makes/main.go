package main

import (
	"encoding/json"
	"log"
	"net/http"
	"valuation-api/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var year string = request.PathParameters["year"]

	data, statusCode, err := call(&year)

	if err != nil {
		log.Fatal(err)
	}

	if statusCode >= 400 {
		return utils.HandleError(err, statusCode)
	}

	d, err := json.Marshal(data)

	if err != nil {
		log.Fatal(err)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(d),
		StatusCode: http.StatusOK,
	}, nil
}

type MakeReturnValue struct {
	MakeID          int    `json:"MakeID"`
	MakeDisplayName string `json:"MakeDisplayName"`
}

type MakeIDCategories struct {
	MakeID       int    `json:"MakeID"`
	CategoryID   int    `json:"CategoryID"`
	CategoryName string `json:"CategoryName"`
}

type Make struct {
	ErrorReturnTO struct {
		ErrorMessage string `json:"ErrorMessage"`
	} `json:"ErrorReturnTO"`
	MakeDisplayName string `json:"MakeDisplayName"`
	MakeNotes       []struct {
		MakeNoteID   int    `json:"MakeNoteID"`
		MakeNoteText string `json:"MakeNoteText"`
	} `json:"MakeNotes"`
	MakeIDCategories []MakeIDCategories `json:"MakeIDCategories"`
	VersionTO        struct {
		VersionID   int    `json:"VersionID"`
		VersionName string `json:"VersionName"`
	} `json:"VersionTO"`
}

type GetMakesResult struct {
	GetMakesResult []Make `json:"GetMakesResult"`
}

func call(year *string) ([]MakeReturnValue, int, error) {

	body, statusCode, err := utils.MakeRequest("GET", "UsedPowersportsService.svc/Makes/"+*year+"/0/1")

	if err != nil {
		return nil, statusCode, err
	}

	var makes *GetMakesResult
	json.Unmarshal([]byte(string(body)), &makes)

	result := transform(&makes.GetMakesResult)
	return result, statusCode, nil
}

func transform(makes *[]Make) []MakeReturnValue {
	var result []MakeReturnValue

	for _, make := range *makes {
		var item MakeReturnValue = MakeReturnValue{
			MakeDisplayName: make.MakeDisplayName,
			MakeID:          make.MakeIDCategories[0].MakeID,
		}
		result = append(result, item)
	}

	return result
}

func main() {
	lambda.Start(Handler)
}
