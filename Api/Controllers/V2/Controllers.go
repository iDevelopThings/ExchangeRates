package V2

import (
	"strconv"

	"Currencies/Api/Controllers"
	"Currencies/Currency"
)

type V2Controllers struct {
	Controllers.ApiServiceControllers
}

type ConvertAmountResponse struct {
	From   string  `json:"from,omitempty"`
	To     string  `json:"to,omitempty"`
	Rate   float64 `json:"rate,omitempty"`
	Amount float64 `json:"amount,omitempty"`
}

type AllRatesResponse map[string]map[string]float64

func (v *V2Controllers) Convert(ctx Controllers.RequestContextImpl) (ConvertAmountResponse, Controllers.ErrorInterface) {
	from := ctx.CurrencyCodeFromUrl("from", Currency.EUR)
	to := ctx.CurrencyCodeFromUrl("to", Currency.EUR)

	converted := ConvertAmountResponse{
		From:   string(from),
		To:     string(to),
		Rate:   0,
		Amount: 0,
	}

	convertedRate, err := ctx.GetRates(to).Convert(from)
	if err != nil {
		return converted, ctx.CreateError("Failed to convert currencies", 400)
	}
	converted.Rate = convertedRate

	amount, amountOk := ctx.Params()["amount"]

	if amountOk {
		parsedAmount, err := strconv.ParseFloat(amount, 64)
		if err != nil {
			return converted, ctx.CreateError("Failed to convert amount to float64", 400)
		}
		convertedAmount, err := ctx.GetRates(to).ConvertAmount(parsedAmount, from)
		if err != nil {
			return converted, ctx.CreateError("Failed to convert amount", 400)
		}

		converted.Amount = convertedAmount
	}

	return converted, nil
}
