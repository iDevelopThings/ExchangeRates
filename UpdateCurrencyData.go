package main

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	
	"github.com/go-redis/redis"
)

type ecbCurrencies struct {
	Currencies []struct {
		Currency string  `xml:"currency,attr"`
		Rate     float32 `xml:"rate,attr"`
	} `xml:"Cube>Cube>Cube"`
}

type ecbUpdateDate struct {
	Date struct {
		Time string `xml:"time,attr"`
	} `xml:"Cube>Cube"`
}

type CurrencyData map[string]Currency

type Currency struct {
	Currency string  `json:"currency"`
	Rate     float32 `json:"rate"`
}

func UpdateCurrencyData() (*ecbCurrencies, *ecbUpdateDate, error) {
	response, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil, nil, err
	}
	
	defer response.Body.Close()
	
	var currencies ecbCurrencies
	var updateDate ecbUpdateDate
	
	data, err := ioutil.ReadAll(response.Body)
	
	if err := xml.Unmarshal(data, &currencies); err != nil {
		return nil, nil, err
	}
	if err := xml.Unmarshal(data, &updateDate); err != nil {
		return nil, nil, err
	}
	return &currencies, &updateDate, nil
}

func StoreCurrencies() (*CurrencyData, error) {
	currencies, _, err := UpdateCurrencyData()
	
	if err != nil {
		return nil, err
	}
	
	currencyData := make(CurrencyData)
	
	for _, currency := range currencies.Currencies {
		currencyData[currency.Currency] = Currency{
			Currency: currency.Currency,
			Rate:     currency.Rate,
		}
	}
	
	err = Set("currencies", currencyData)
	
	if err != nil {
		return nil, err
	}
	
	currencyData = CurrencyData{}
	
	err = Get("currencies", &currencyData)
	
	if err != nil {
		return nil, err
	}
	
	return &currencyData, nil
}

func GetOrUpdateCurrencies() (*CurrencyData, error) {
	var currencies CurrencyData
	
	key := "currencies"
	
	if err := Get(key, &currencies); err != nil {
		if err == redis.Nil {
			currencies, err := StoreCurrencies()
			
			return currencies, err
		}
	}
	
	return &currencies, nil
	
}