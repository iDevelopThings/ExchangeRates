package V1

import (
	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

func (v *V1Controllers) SingleConversion(ctx Controllers.RequestContextImpl) {
	from := ctx.CurrencyCodeFromUrl("from", Currency.EUR)
	to := ctx.CurrencyCodeFromUrl("to", Currency.EUR)

	convertedRate, err := ctx.GetRates(to).Convert(from)
	if err != nil {
		ctx.Error("Failed to convert currencies", 400)
		return
	}

	ctx.Json(SingleConversionsResponse{
		Base: from.ToString(),
		Currency: CurrencyConversion{
			From: from.ToString(),
			To:   to.ToString(),
			Rate: convertedRate,
		},
	})
}
