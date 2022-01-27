package Currency

import (
	"math"
)

func ConvertCurrency(fromRate, toRate float64) float64 {
	return Round64(Round64(toRate, 4)/Round64(fromRate, 4), 4)
}
func ConvertCurrencyAmount(value float64, fromRate, toRate float64) float64 {
	return Round64(value*Round64(toRate, 4)/Round64(fromRate, 4), 4)
}

func Round64(x float64, prec int) float64 {
	if math.IsNaN(x) || math.IsInf(x, 0) {
		return x
	}

	sign := 1.0
	if x < 0 {
		sign = -1
		x *= -1
	}

	var rounder float64
	pow := math.Pow(10, float64(prec))
	intermed := x * pow
	_, frac := math.Modf(intermed)

	if frac >= 0.5 {
		rounder = math.Ceil(intermed)
	} else {
		rounder = math.Floor(intermed)
	}

	return rounder / pow * sign
}
