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
	vin := request.PathParameters["vin"]

	data, statusCode, err := call(&vin)

	if err != nil {
		log.Fatal(err)
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

type Vin struct {
	Make             string `json:"Make"`
	MakeID           string `json:"MakeID"`
	Year             int    `json:"Year"`
	Model            string `json:"Model"`
	ModelID          int    `json:"ModelID"`
	ModelType        string `json:"ModelType"`
	ModelNo          string `json:"ModelNo"`
	MSRP             string `json:"MSRP"`
	LowRetail        string `json:"LowRetail"`
	AverageRetail    int    `json:"AverageRetail"`
	LowTrade         string `json:"LowTrade"`
	HighTrade        string `json:"HighTrade"`
	AverageWholesale string `json:"AverageWholesale"`
	Transmission     string `json:"Transmission"`
	Weight           string `json:"Weight"`
	EngineCC         string `json:"EngineCC"`
	Stroke           string `json:"Stroke"`
	Cylinders        string `json:"Cylinders"`
}

type ModelsByVINResult struct {
	Status int   `json:"Status"`
	Models []Vin `json:"Models"`
}

type GetModelsByVINResult struct {
	GetModelsByVINResult ModelsByVINResult `json:"GetModelsByVINResult"`
}

func call(vin *string) (utils.PricesReturnValue, int, error) {

	vinBody, statusCode, err := utils.MakeRequest("GET", "UsedPowersportsService.svc/VIN/"+*vin)

	var vinResult *GetModelsByVINResult
	json.Unmarshal([]byte(string(vinBody)), &vinResult)

	tpBody, tpStatusCode, tpErr := utils.PriceCall(&vinResult.GetModelsByVINResult.Models[0].ModelID)

	if tpErr != nil {
		return utils.PricesReturnValue{}, tpStatusCode, tpErr
	}

	return tpBody, statusCode, err
}

func main() {
	lambda.Start(Handler)
}
