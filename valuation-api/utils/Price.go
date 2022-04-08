package utils

import (
	"encoding/json"
	"log"
	"strconv"
)

type PricesReturnValue struct {
	TotalPrice    float32 `json:"TotalPrice"`
	AverageRetail float32 `json:"AverageRetail"`
	CleanTrade    float32 `json:"CleanTrade"`
	RoughTrade    float32 `json:"RoughTrade"`
	AverageTrade  float32 `json:"AverageTrade"`
	Make          string  `json:"Make"`
	Model         string  `json:"Model"`
	Year          string  `json:"Year"`
}

type TradeValueReport struct {
	GetTradePricesResult GetTradePricesResult `json:"GetTradePricesResult"`
}

type GetTradePricesResult struct {
	Make                 string     `json:"Make"`
	Model                string     `json:"Model"`
	ModelID              string     `json:"ModelID"`
	ModelType            string     `json:"ModelType"`
	Year                 string     `json:"Year"`
	Category             string     `json:"Category"`
	CategoryID           int        `json:"CategoryID"`
	VersionID            int        `json:"VersionID"`
	VersionName          string     `json:"VersionName"`
	BasePrices           PricesType `json:"BasePrices"`
	TotalPrices          PricesType `json:"TotalPrices"`
	PriceTypeDefinitions []string
	SelectedOptions      []SelectedOptions `json:"SelectedOptions"`
	ReportedEngine       struct {
		DefaultEngine         bool   `json:"DefaultEngine"`
		EngineId              int    `json:"EngineId"`
		EngineName            string `json:"EngineName"`
		EngineNote            string `json:"EngineNote"`
		EnginePriceAdjustment int    `json:"EnginePriceAdjustment"`
	} `json:"ReportedEngine"`
	EnginePriceAdjustment int `json:"EnginePriceAdjustment"`
	Specs                 struct {
		// __type       string `json:"__type"`
		CC           int    `json:"CC"`
		Cylinders    int    `json:"Cylinders"`
		Stroke       int    `json:"Stroke"`
		Transmission string `json:"Transmission"`
		Weight       string `json:"Weight"`
	} `json:"Specs"`
}

type PricesType struct {
	// __type                string `json:"__type"`
	AverageRetail         float32 `json:"AverageRetail"`
	CleanTradeInWholesale float32 `json:"CleanTradeInWholesale"`
	RoughTradeInWholesale float32 `json:"RoughTradeInWholesale"`
	SuggestedListPrice    float32 `json:"SuggestedListPrice"`
}

type SelectedOptions struct {
	OptionGroup       string       `json:"OptionGroup"`
	OptionGroupID     int          `json:"OptionGroupID"`
	OptionDisplayName string       `json:"OptionDisplayName"`
	OptionCode        string       `json:"OptionCode"`
	Values            OptionValues `json:"Values"`
}

type OptionValues struct {
	// __type        string `json:"__type"`
	AverageRetail float32 `json:"AverageRetail"`
	CleanTrade    float32 `json:"CleanTrade"`
	RoughTrade    float32 `json:"RoughTrade"`
}

func PriceCall(modelID *int) (PricesReturnValue, int, error) {
	model := strconv.Itoa(*modelID)

	body, statusCode, err := MakeRequest("GET", "UsedPowersportsService.svc/TradeValueReport/"+model+"/0/0")

	if err != nil {
		log.Fatal(err)
		return PricesReturnValue{}, statusCode, err
	}

	var price *TradeValueReport
	json.Unmarshal([]byte(string(body)), &price)

	result := transform(&price.GetTradePricesResult)
	return result, statusCode, nil
}

func transform(price *GetTradePricesResult) PricesReturnValue {
	var result PricesReturnValue = PricesReturnValue{
		TotalPrice:    (*price).TotalPrices.CleanTradeInWholesale,
		AverageRetail: (*price).TotalPrices.AverageRetail,
		CleanTrade:    (*price).TotalPrices.CleanTradeInWholesale,
		RoughTrade:    (*price).TotalPrices.RoughTradeInWholesale,
		AverageTrade:  ((*price).TotalPrices.CleanTradeInWholesale + (*price).TotalPrices.RoughTradeInWholesale) / 2,
		Year:          (*price).Year,
		Model:         (*price).Model,
		Make:          (*price).Make,
	}

	return result
}
