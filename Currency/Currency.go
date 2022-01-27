package Currency

type CurrencyCode string

func (c CurrencyCode) IsValid() bool {
	_, isOk := ValidCurrencies[c]

	return isOk
}

func (c CurrencyCode) ToString() string {
	return string(c)
}
