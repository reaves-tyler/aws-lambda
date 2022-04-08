package main

import (
	"encoding/json"
	"net/http"
	"valuation-api/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	make := request.PathParameters["make"]
	year := request.PathParameters["year"]

	data, statusCode, err := call(&make, &year)

	if err != nil {
		return utils.HandleError(err, statusCode)
	}

	if statusCode >= 400 {
		return utils.HandleError(err, statusCode)
	}

	d, err := json.Marshal(data)

	if err != nil {
		return utils.HandleError(err, http.StatusInternalServerError)
	}

	return events.APIGatewayProxyResponse{
		Body:       string(d),
		StatusCode: http.StatusOK,
	}, nil

}

type ModelsReturnValue struct {
	ModelID   int    `json:"ModelID"`
	ModelName string `json:"ModelName"`
}

type ModelTrims struct {
	ModelName    string `json:"ModelName"`
	ModelID      int    `json:"ModelID"`
	CategoryID   int    `json:"CategoryID"`
	CategoryName string `json:"CategoryName"`
	ModelTypeID  int    `json:"ModelTypeID"`
	ModelType    string `json:"ModelType"`
}

type ModelTrimsResult struct {
	Make struct {
		MakeID          int    `json:"MakeID"`
		MakeDisplayName string `json:"MakeDisplayName"`
	} `json:"Make"`
	MakeNotes  []string     `json:"MakeNotes"`
	ModelTrims []ModelTrims `json:"ModelTrims"`
}

type GetModelTrimsResult struct {
	GetModelTrimsResult ModelTrimsResult `json:"GetModelTrimsResult"`
}

func call(make *string, year *string) ([]ModelsReturnValue, int, error) {

	body, statusCode, err := utils.MakeRequest("GET", "UsedPowersportsService.svc/Models/"+*make+"/"+*year+"/0")

	if err != nil {
		return nil, statusCode, err
	}

	var models *GetModelTrimsResult
	json.Unmarshal([]byte(string(body)), &models)

	result := transform(&models.GetModelTrimsResult.ModelTrims)
	return result, statusCode, nil
}

func transform(models *[]ModelTrims) []ModelsReturnValue {
	var result []ModelsReturnValue

	for i := range *models {
		var item ModelsReturnValue = ModelsReturnValue{
			ModelID:   (*models)[i].ModelID,
			ModelName: (*models)[i].ModelName,
		}

		result = append(result, item)
	}
	return result
}

func main() {
	lambda.Start(Handler)
}
