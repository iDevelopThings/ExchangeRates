package main

import (
	"fmt"

	"Currencies/Api"
	"Currencies/Currency"
	"github.com/joho/godotenv"
)

var api *Api.ApiService
var currencyRates *Currency.CurrencyDataRates

func main() {
	var err error

	err = godotenv.Load()
	if err != nil {
		fmt.Println(err)
	}

	currencyRates, err = Currency.NewCurrencyDataProcessor()
	if err != nil {
		fmt.Println(err)
	}

	api = Api.NewApiService(currencyRates)
	api.Listen()
}
