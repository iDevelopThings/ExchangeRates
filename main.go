package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type CurrencyConversion struct {
	From     string  `json:"from"`
	To       string  `json:"to"`
	Rate     float32 `json:"rate"`
	Currency string  `json:"currency"`
}

type ConversionsResponse struct {
	Base       string                        `json:"base"`
	Currencies map[string]CurrencyConversion `json:"currencies"`
}

func main() {
	godotenv.Load()

	InitRedis()

	router := mux.NewRouter()

	router.HandleFunc("/latest", CORS(Latest))

	fmt.Println("Running on http://127.0.0.1:" + os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), router)
}

func CORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			w.Header().Set("Access-Control-Allow-Headers",
				"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		}
		// Preflighted OPTIONS request
		if r.Method == "OPTIONS" {
			return
		}
		fn(w, r)
	}
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
			From:     baseCurrency.Currency,
			To:       currency.Currency,
			Rate:     currency.Rate,
			Currency: currency.Currency,
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
