package V1

import (
	"net/http"

	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

func (v *V1Controllers) ListAllConversions(ctx Controllers.RequestContextImpl) {
	baseCurrency := ctx.CurrencyCodeFromUrl("base", Currency.EUR)

	if !baseCurrency.IsValid() {
		ctx.Error("Invalid currency passed.", http.StatusBadRequest)
		return
	}

	responseData := &ConversionsResponse{
		Base:       baseCurrency,
		Currencies: make(map[string]CurrencyConversion),
	}

	for currencyCode, rate := range ctx.GetRates(baseCurrency).Rates {
		responseData.Currencies[string(currencyCode)] = CurrencyConversion{
			From: string(baseCurrency),
			To:   string(currencyCode),
			Rate: rate,
		}
	}

	ctx.Json(responseData)
}
