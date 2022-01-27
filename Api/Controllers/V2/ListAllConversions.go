package V2

import (
	"net/http"

	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

func (v *V2Controllers) ListAllConversions(ctx Controllers.RequestContextImpl) {
	baseCurrency := ctx.CurrencyCodeFromUrl("base", Currency.EUR)

	if !baseCurrency.IsValid() {
		ctx.Error("Invalid currency passed.", http.StatusBadRequest)
		return
	}

	ctx.Json(ctx.GetRates(baseCurrency).Rates)
}

func (v *V2Controllers) ListAll(ctx Controllers.RequestContextImpl) {
	response := make(AllRatesResponse)

	for code, rate := range ctx.Rates().Rates {
		for currencyCode, conversion := range rate.Rates {
			_, added := response[code.ToString()]
			if !added {
				response[code.ToString()] = make(map[string]float64)
			}
			response[code.ToString()][currencyCode.ToString()] = conversion
		}
	}

	ctx.Json(response)
}
