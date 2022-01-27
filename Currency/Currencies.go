package Currency

const (
	AUD CurrencyCode = "AUD" // Australian Dollar (A$)
	BGN CurrencyCode = "BGN" // Bulgarian Lev (BGN)
	BRL CurrencyCode = "BRL" // Brazilian Real (R$)
	CAD CurrencyCode = "CAD" // Canadian Dollar (CA$)
	CHF CurrencyCode = "CHF" // Swiss Franc (CHF)
	CNY CurrencyCode = "CNY" // Chinese Yuan (CN¥)
	CZK CurrencyCode = "CZK" // Czech Republic Koruna (CZK)
	DKK CurrencyCode = "DKK" // Danish Krone (DKK)
	EUR CurrencyCode = "EUR" // Euro (€)
	GBP CurrencyCode = "GBP" // British Pound Sterling (£)
	HKD CurrencyCode = "HKD" // Hong Kong Dollar (HK$)
	HRK CurrencyCode = "HRK" // Croatian Kuna (HRK)
	HUF CurrencyCode = "HUF" // Hungarian Forint (HUF)
	IDR CurrencyCode = "IDR" // Indonesian Rupiah (IDR)
	ILS CurrencyCode = "ILS" // Israeli New Sheqel (₪)
	INR CurrencyCode = "INR" // Indian Rupee (Rs.)
	JPY CurrencyCode = "JPY" // Japanese Yen (¥)
	KRW CurrencyCode = "KRW" // South Korean Won (₩)
	LTL CurrencyCode = "LTL" // Lithuanian Litas (LTL)
	MXN CurrencyCode = "MXN" // Mexican Peso (MX$)
	MYR CurrencyCode = "MYR" // Malaysian Ringgit (MYR)
	NOK CurrencyCode = "NOK" // Norwegian Krone (NOK)
	NZD CurrencyCode = "NZD" // New Zealand Dollar (NZ$)
	PHP CurrencyCode = "PHP" // Philippine Peso (Php)
	PLN CurrencyCode = "PLN" // Polish Zloty (PLN)
	RON CurrencyCode = "RON" // Romanian Leu (RON)
	RUB CurrencyCode = "RUB" // Russian Ruble (RUB)
	SEK CurrencyCode = "SEK" // Swedish Krona (SEK)
	SGD CurrencyCode = "SGD" // Singapore Dollar (SGD)
	THB CurrencyCode = "THB" // Thai Baht (฿)
	TRY CurrencyCode = "TRY" // Turkish Lira (TRY)
	USD CurrencyCode = "USD" // US Dollar ($)
	ZAR CurrencyCode = "ZAR" // South African Rand (ZAR)
)

var ValidCurrencies = map[CurrencyCode]bool{
	AUD: true,
	BGN: true,
	BRL: true,
	CAD: true,
	CHF: true,
	CNY: true,
	CZK: true,
	DKK: true,
	EUR: true,
	GBP: true,
	HKD: true,
	HRK: true,
	HUF: true,
	IDR: true,
	ILS: true,
	INR: true,
	JPY: true,
	KRW: true,
	LTL: true,
	MXN: true,
	MYR: true,
	NOK: true,
	NZD: true,
	PHP: true,
	PLN: true,
	RON: true,
	RUB: true,
	SEK: true,
	SGD: true,
	THB: true,
	TRY: true,
	USD: true,
	ZAR: true,
}
