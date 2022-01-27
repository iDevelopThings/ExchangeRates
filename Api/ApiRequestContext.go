package Api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"Currencies/Api/Controllers"
	"Currencies/Currency"
	"github.com/naoina/denco"
)

type RequestContext struct {
	Request  *http.Request
	Response http.ResponseWriter
	params   map[string]string
	service  *ApiService
}

type ErrorResponse struct {
	code    int    `json:"error,omitempty"`
	message string `json:"message,omitempty"`
	ctx     *RequestContext
}

func (e *ErrorResponse) Error() string {
	return "Error: " + string(rune(e.code)) + "; " + e.message
}

func (e *ErrorResponse) Message() string {
	return e.message
}

func (e *ErrorResponse) Code() int {
	return e.code
}

func (e *ErrorResponse) Response(ctx *RequestContext) Controllers.ErrorInterface {
	e.ctx = ctx
	return e
}

func (e *ErrorResponse) SendError() {
	if e.ctx == nil {
		log.Fatal("Trying to send ErrorResponse, but ctx is not set.")
		return
	}

	e.ctx.Error(e.message, e.code)
}

func NewError(message string, code int) *ErrorResponse {
	return &ErrorResponse{
		code:    code,
		message: message,
	}
}

func NewRequestContext(request *http.Request, response http.ResponseWriter, params denco.Params, service *ApiService) *RequestContext {
	mappedParams := map[string]string{}

	for _, param := range params {
		if param.Value != "" {
			mappedParams[param.Name] = param.Value
		}
	}

	return &RequestContext{
		Request:  request,
		Response: response,
		service:  service,
		params:   mappedParams,
	}
}

func (ctx *RequestContext) Json(data interface{}) {

	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.Header().Set("X-Content-Type-Options", "nosniff")

	response, err := json.Marshal(data)
	if err != nil {
		http.Error(ctx.Response, http.StatusText(500), 500)
		log.Fatalf("Error encoding JSON: %s", err)
		return
	}

	_, err = ctx.Response.Write(response)
	if err != nil {
		log.Fatalf("Error encoding JSON: %s", err)
		return
	}
}

func (ctx *RequestContext) CreateError(message string, code int) Controllers.ErrorInterface {
	return NewError(message, code).Response(ctx)
}

func (ctx *RequestContext) Error(message string, code int) {
	log.Printf("Error response: %s \n", message)

	ctx.Response.Header().Set("Content-Type", "application/json; charset=utf-8")
	ctx.Response.Header().Set("X-Content-Type-Options", "nosniff")

	ctx.Response.WriteHeader(code)

	ctx.Json(NewError(message, code))
}

func (ctx *RequestContext) Html(fileName string) {
	ctx.service.HtmlRenderer.HTML(ctx.Response, http.StatusOK, fileName, nil)
}

func (ctx *RequestContext) UrlParam(key string, defaultValue ...interface{}) string {
	if baseParam := ctx.Request.URL.Query().Get(key); baseParam != "" {
		return baseParam
	}

	if defaultValue == nil || len(defaultValue) == 0 {
		return ""
	}

	strVal := defaultValue[0]

	switch v := strVal.(type) {
	case Currency.CurrencyCode:
		return v.ToString()
	case string:
		return fmt.Sprintf("%v", v)
	}

	return fmt.Sprintf("%v", defaultValue[0])
}

func (ctx *RequestContext) UrlParamUpperCase(key string, defaultValue ...interface{}) string {
	return strings.ToUpper(ctx.UrlParam(key, defaultValue))
}

func (ctx *RequestContext) CurrencyCodeFromUrl(key string, defaultValue Currency.CurrencyCode) Currency.CurrencyCode {
	baseParam := ctx.Request.URL.Query().Get(key)
	if baseParam == "" {
		if !ctx.HasVar(key) {
			return defaultValue
		}

		param, hasParam := ctx.params[key]
		if !hasParam {
			return defaultValue
		}

		baseParam = param
	}

	currency := Currency.CurrencyCode(strings.ToUpper(baseParam))

	if !currency.IsValid() {
		ctx.Error(fmt.Sprintf("You specified an invalid currency for \"%s\"", key), http.StatusBadRequest)
	}

	return currency
}

func (ctx *RequestContext) HasVar(key string) bool {
	_, isOk := ctx.params[key]

	return isOk
}

func (ctx *RequestContext) GetRates(rate Currency.CurrencyCode) *Currency.Rate {
	return ctx.service.CurrencyRates.GetRates(rate)
}

func (ctx *RequestContext) Rates() *Currency.CurrencyDataRates {
	return ctx.service.CurrencyRates
}

func (ctx *RequestContext) Params() map[string]string {
	return ctx.params
}

func (ctx *RequestContext) Get() Controllers.RequestContextImpl {
	return ctx
}
