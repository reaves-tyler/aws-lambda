package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"valuation-api/utils"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	modelID := request.PathParameters["modelID"]

	modID, err := strconv.Atoi(modelID)

	if err != nil {
		return utils.HandleError(err, http.StatusInternalServerError)
	}

	data, statusCode, err := utils.PriceCall(&modID)

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

func main() {
	lambda.Start(Handler)
}
