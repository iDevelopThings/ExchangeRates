package Controllers

import "Currencies/Currency"

type ErrorInterface interface {
	Message() string
	Code() int
	Error() string
	SendError()
}

type RequestContextImpl interface {
	Params() map[string]string
	Json(data interface{})
	CreateError(message string, code int) ErrorInterface
	Error(message string, code int)
	Html(fileName string)
	UrlParam(key string, defaultValue ...interface{}) string
	UrlParamUpperCase(key string, defaultValue ...interface{}) string
	CurrencyCodeFromUrl(key string, defaultValue Currency.CurrencyCode) Currency.CurrencyCode
	HasVar(key string) bool
	GetRates(rate Currency.CurrencyCode) *Currency.Rate
	Rates() *Currency.CurrencyDataRates
}
