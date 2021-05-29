package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/thedevsaddam/renderer"
)

type CurrencyConversion struct {
	From string  `json:"from"`
	To   string  `json:"to"`
	Rate float64 `json:"rate"`
}

type ConversionsResponse struct {
	Base       string                        `json:"base"`
	Currencies map[string]CurrencyConversion `json:"currencies"`
}
type SingleConversionsResponse struct {
	Base     string             `json:"base"`
	Currency CurrencyConversion `json:"currency"`
}
type ConvertAmountResponse struct {
	Base     string             `json:"base"`
	Currency CurrencyConversion `json:"currency"`
	Amount   float64            `json:"amount"`
}

var rnd *renderer.Render

func main() {
	godotenv.Load()
	opts := renderer.Options{
		ParseGlobPattern: "./*.html",
	}

	rnd = renderer.New(opts)

	InitRedis()

	router := mux.NewRouter()

	router.HandleFunc("/", Information)
	router.HandleFunc("/latest", CORS(Latest))
	router.HandleFunc("/conversion/{from}/{to}", CORS(Conversion))
	router.HandleFunc("/convert/{from}/{to}/{amount}", CORS(ConvertAmount))

	fmt.Println("Running on http://127.0.0.1:" + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

func CORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}

		if r.Method == "OPTIONS" {
			return
		}

		fn(w, r)
	}
}

func Information(w http.ResponseWriter, r *http.Request) {
	rnd.HTML(w, http.StatusOK, "info", nil)
}

func ConvertAmount(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	from, fromOk := params["from"]
	to, toOk := params["to"]
	amount, amountOk := params["amount"]

	if !fromOk || !toOk || !amountOk {
		http.Error(w, "Missing from, to or amount in request url. Should be: /convert/usd/gbp/1.50", 400)
		return
	}

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	_, err := GetOrUpdateCurrencies()

	if err != nil {
		http.Error(w, "Something went wrong getting currency information...", 500)
		return
	}

	currencies := map[string]CurrencyConversion{}
	err = Get("conversions-"+from, &currencies)

	if err != nil {
		http.Error(w, "Something went wrong getting currency information..."+to, 500)
		return
	}

	toCurrency, toCurrencyOk := currencies[to]

	if !toCurrencyOk {
		http.Error(w, "Cannot get currency info for "+to+". Is this a valid currency?", 400)
		return
	}

	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		http.Error(w, "Failed to convert amount to float64", 400)
		return
	}

	responseData := &ConvertAmountResponse{
		Base:     from,
		Currency: toCurrency,
		Amount:   parsedAmount * toCurrency.Rate,
	}

	response, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func Conversion(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	from, fromOk := params["from"]
	to, toOk := params["to"]

	if !fromOk || !toOk {
		http.Error(w, "Missing from or to in request url. Should be: /conversion/usd/gbp", 400)
		return
	}

	from = strings.ToUpper(from)
	to = strings.ToUpper(to)

	_, err := GetOrUpdateCurrencies()

	if err != nil {
		http.Error(w, "Something went wrong getting currency information...", 500)
		return
	}

	currencies := map[string]CurrencyConversion{}
	err = Get("conversions-"+from, &currencies)

	if err != nil {
		http.Error(w, "Something went wrong getting currency information..."+to, 500)
		return
	}

	toCurrency, toCurrencyOk := currencies[to]

	if !toCurrencyOk {
		http.Error(w, "Cannot get currency info for "+to+". Is this a valid currency?", 400)
		return
	}

	responseData := &SingleConversionsResponse{
		Base:     from,
		Currency: toCurrency,
	}

	response, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func Latest(w http.ResponseWriter, r *http.Request) {
	baseParam := r.URL.Query().Get("base")

	if baseParam == "" {
		baseParam = "EUR"
	}
	responseData := &ConversionsResponse{
		Base:       baseParam,
		Currencies: make(map[string]CurrencyConversion),
	}

	currenciesMap, err := GetOrUpdateCurrencies()

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	currencies := *currenciesMap

	for _, currency := range currencies {
		baseCurrency := currencies[strings.ToUpper(baseParam)]

		if baseParam != "EUR" {
			currency.Rate = currency.Rate / baseCurrency.Rate
		}

		responseData.Currencies[currency.Currency] = CurrencyConversion{
			From: baseCurrency.Currency,
			To:   currency.Currency,
			Rate: float64(currency.Rate),
		}
	}

	response, err := json.Marshal(responseData)

	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}
