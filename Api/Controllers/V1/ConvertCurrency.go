package V1

import (
	"strconv"

	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

func (v *V1Controllers) ConvertCurrency(ctx Controllers.RequestContextImpl) {
	from := ctx.CurrencyCodeFromUrl("from", Currency.EUR)
	to := ctx.CurrencyCodeFromUrl("to", Currency.EUR)

	amount, amountOk := ctx.Params()["amount"]

	if !amountOk {
		ctx.Error("Missing from, to or amount in request url. Should be: /convert/usd/gbp/1.50", 400)
		return
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		ctx.Error("Failed to convert amount to float64", 400)
		return
	}

	convertedRate, err := ctx.GetRates(to).Convert(from)
	convertedAmount, err := ctx.GetRates(to).ConvertAmount(parsedAmount, from)
	if err != nil {
		ctx.Error("Failed to convert currencies", 400)
		return
	}

	ctx.Json(ConvertAmountResponse{
		Base: string(from),
		Currency: CurrencyConversion{
			From: string(from),
			To:   string(to),
			Rate: convertedRate,
		},
		Amount: convertedAmount,
	})
}
