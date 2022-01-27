package Currency

import (
	"errors"
)

type Rate struct {
	Rates    map[CurrencyCode]float64
	Currency CurrencyCode
}

func (r *Rate) ConvertAmount(value float64, to CurrencyCode) (float64, error) {
	if r.Rates[to] == 0 {
		return 0, errors.New("invalid currency provided")
	}

	return ConvertCurrencyAmount(value, r.Rates[to], r.Rates[r.Currency]), nil
}

func (r *Rate) Convert(to CurrencyCode) (float64, error) {
	if r.Rates[to] == 0 {
		return 0, errors.New("invalid currency provided")
	}

	return ConvertCurrency(r.Rates[to], r.Rates[r.Currency]), nil
}
