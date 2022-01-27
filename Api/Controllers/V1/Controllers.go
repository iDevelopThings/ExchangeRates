package V1

import (
	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

type V1Controllers struct {
	Controllers.ApiServiceControllers
}

type CurrencyConversion struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type ConvertAmountResponse struct {
	Base     string             `json:"base"`
	Currency CurrencyConversion `json:"currency"`
	Amount   float64            `json:"amount"`
}

type ConversionsResponse struct {
	Base       Currency.CurrencyCode         `json:"base"`
	Currencies map[string]CurrencyConversion `json:"currencies"`
}

type SingleConversionsResponse struct {
	Base     string             `json:"base"`
	Currency CurrencyConversion `json:"currency"`
}
