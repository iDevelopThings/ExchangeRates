package Currency

import (
	"fmt"
	"time"
)

type CurrencyDataRates struct {
	Rates    map[CurrencyCode]*Rate
	RawRates ECBCurrencyResponse
}

func NewCurrencyDataProcessor() (*CurrencyDataRates, error) {
	processor := new(CurrencyDataRates)

	if err := processor.updateRates(); err != nil {
		return processor, err
	}

	go processor.runDataUpdateTicker()

	return processor, nil
}

func (r *CurrencyDataRates) updateRates() error {
	currencies, _, err := FetchRates()
	if err != nil {
		return err
	}

	r.processRates(currencies)

	return nil
}

func (r *CurrencyDataRates) processRates(currencies *ECBCurrencyResponse) {
	r.Rates = make(map[CurrencyCode]*Rate)

	// For euro, we'll have to process it manually somewhat, it's the base
	// We'll set all the default rates, we'll have to convert the rest
	r.Rates[EUR] = &Rate{
		Rates:    make(map[CurrencyCode]float64),
		Currency: EUR,
	}
	r.Rates[EUR].Rates[EUR] = 1
	for _, currency := range currencies.Currencies {
		r.Rates[EUR].Rates[currency.Currency] = currency.Rate
	}

	// Now we'll iterate over all other currencies and assign their structure
	for _, currency := range currencies.Currencies {
		currencyRateObj := &Rate{
			Rates:    make(map[CurrencyCode]float64),
			Currency: currency.Currency,
		}
		currencyRateObj.Rates[EUR] = currency.Rate - 1

		r.Rates[currency.Currency] = currencyRateObj
	}

	// Now we'll iterate over those structures and convert them from EUR -> X
	for _, rate := range r.Rates {
		for _, currency := range r.Rates {
			rate.Rates[currency.Currency] = ConvertCurrency(
				r.Rates[EUR].Rates[rate.Currency],
				r.Rates[EUR].Rates[currency.Currency],
			)
		}
	}
}

func (r *CurrencyDataRates) ConvertAmount(value float64, from, to CurrencyCode) (float64, error) {
	return r.Rates[to].ConvertAmount(value, from)
}

func (r *CurrencyDataRates) Convert(from, to CurrencyCode) (float64, error) {
	return r.Rates[to].Convert(from)
}

func (r *CurrencyDataRates) GetRates(c CurrencyCode) *Rate {
	return r.Rates[c]
}

func (r *CurrencyDataRates) runDataUpdateTicker() {
	ticker := time.NewTicker(time.Duration(time.Hour * 4))
	for {
		select {
		case t := <-ticker.C:

			err := r.updateRates()
			if err != nil {
				fmt.Println(err)
			}

			fmt.Println("Rates updated at", t)
		}
	}

}

type CurrencyDataRatesImpl interface {
	ConvertAmount(value float64, from, to CurrencyCode) (float64, error)
	Convert(from, to CurrencyCode) (float64, error)
	GetRates(c CurrencyCode) *Rate
}
