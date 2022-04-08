package utils

import (
	"encoding/json"
	"os"
)

type Token struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func GetToken() (string, int, error) {
	body, statusCode, err := MakeTokenRequest()

	if err != nil {
		return "", statusCode, err
	}

	var token *Token
	json.Unmarshal([]byte(string(body)), &token)

	return token.AccessToken, statusCode, nil
}

type CustomerByEmail struct {
	Pagination struct {
		Offset int `json:"offset"`
		Limit  int `json:"limit"`
		Total  int `json:"total"`
	} `json:"pagination"`
	Result []struct {
		URL    string `json:"url"`
		ID     string `json:"id"`
		Status string `json:"status"`
		Type   string `json:"type"`
	} `json:"result"`
}

func CheckCustomerExists(email string) (string, int, error) {
	token, statusCode, err := GetToken()
	body, statusCode, err := MakeRequest("GET", os.Getenv("NEBULOUS_BASEURL")+"customers/?email="+email, nil, token)

	if err != nil {
		return "", statusCode, err
	}

	var CustomerByEmail *CustomerByEmail
	json.Unmarshal([]byte(string(body)), &CustomerByEmail)

	if len(CustomerByEmail.Result) == 1 {
		return CustomerByEmail.Result[0].ID, statusCode, nil
	}

	return "", statusCode, nil
}

type CustomerCreate struct {
	Customer CustomerFields `json:"customer"`
}

type CustomerFields struct {
	Email1    string `json:"email1"`
	ShowPhone bool   `json:"showPhone"`
	Type      string `json:"type"`
}

type NewCustomer struct {
	Url    string `json:"url"`
	Result struct {
		ID         string `json:"id"`
		Username   string `json:"username"`
		FirstName  string `json:"firstName"`
		LastName   string `json:"lastName"`
		Name       string `json:"name"`
		Address1   string `json:"address1"`
		Address2   string `json:"address2"`
		City       string `json:"city"`
		StateCode  string `json:"stateCode"`
		Zip        string `json:"zip"`
		CountryId  string `json:"countryId"`
		Email1     string `json:"email1"`
		Email2     string `json:"email2"`
		Status     string `json:"status"`
		Website    string `json:"website"`
		Phone      string `json:"phone"`
		ShowPhone  bool   `json:"showPhone"`
		Notifiable bool   `json:"notifiable"`
		Type       string `json:"type"`
		CreateDate string `json:"createDate"`
	} `json:"result"`
}

func CreateCustomer(email string) (NewCustomer, int, error) {
	token, statusCode, err := GetToken()
	var customer CustomerCreate = CustomerCreate{
		Customer: CustomerFields{
			Email1:    email,
			ShowPhone: false,
			Type:      "customer",
		},
	}

	postData, err := json.Marshal(customer)

	if err != nil {
		return NewCustomer{}, 500, err
	}

	body, statusCode, err := MakeRequest("POST", os.Getenv("NEBULOUS_URL")+"customers/", postData, token)

	if err != nil {
		return NewCustomer{}, statusCode, err
	}

	var newCustomer *NewCustomer
	json.Unmarshal([]byte(string(body)), &newCustomer)

	return *newCustomer, statusCode, nil
}
