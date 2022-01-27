package Currency

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type ECBCurrencyResponse struct {
	Currencies []struct {
		Currency CurrencyCode `xml:"currency,attr"`
		Rate     float64      `xml:"rate,attr"`
	} `xml:"Cube>Cube>Cube"`
}

type ECBUpdateDate struct {
	Date struct {
		Time string `xml:"time,attr"`
	} `xml:"Cube>Cube"`
}

func FetchRates() (*ECBCurrencyResponse, *ECBUpdateDate, error) {
	response, err := http.Get("http://www.ecb.europa.eu/stats/eurofxref/eurofxref-daily.xml")
	if err != nil {
		return nil, nil, err
	}

	defer response.Body.Close()

	var currencies ECBCurrencyResponse
	var updateDate ECBUpdateDate

	data, err := ioutil.ReadAll(response.Body)

	if err := xml.Unmarshal(data, &currencies); err != nil {
		return nil, nil, err
	}

	if err := xml.Unmarshal(data, &updateDate); err != nil {
		return nil, nil, err
	}

	return &currencies, &updateDate, nil
}
