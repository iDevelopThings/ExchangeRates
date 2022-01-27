package Api

import (
	"fmt"
	"net/http"
	"os"

	"Currencies/Api/Controllers"
	"Currencies/Api/Controllers/Docs"
	"Currencies/Api/Controllers/V1"
	"Currencies/Api/Controllers/V2"
	"Currencies/Api/Middleware"
	"Currencies/Currency"
	"github.com/naoina/denco"
	"github.com/thedevsaddam/renderer"
)

type ControllerHandler func(ctx Controllers.RequestContextImpl)
type RouteMapping map[string]denco.HandlerFunc

type ApiService struct {
	CurrencyRates *Currency.CurrencyDataRates
	HtmlRenderer  *renderer.Render

	server       *RequestHandler
	routeMapping RouteMapping
	V1           *V1.V1Controllers
	V2           *V2.V2Controllers
}

func NewApiService(currencyRates *Currency.CurrencyDataRates) *ApiService {

	service := &ApiService{
		CurrencyRates: currencyRates,
		HtmlRenderer:  renderer.New(renderer.Options{ParseGlobPattern: "./*.html"}),
		server:        NewServer(),
		routeMapping:  make(RouteMapping),
		V1:            new(V1.V1Controllers),
		V2:            new(V2.V2Controllers),
	}

	service.addRoute("/", Docs.DocsPage)

	service.mapRoutesForVersion("", service.V1)
	service.mapRoutesForVersion("/v1", service.V1)

	service.mapRoutesForVersion("/v2", service.V2)
	service.addRoute("/v2/all", service.V2.ListAll)

	return service
}

func (service *ApiService) mapRoutesForVersion(prefix string, apiVersion Controllers.ApiServiceControllers) {
	service.addRoute(prefix+"/latest", apiVersion.ListAllConversions)
	service.addRoute(prefix+"/latest/:base", apiVersion.ListAllConversions)

	service.addRoute(prefix+"/currencies", apiVersion.ListAllConversions)
	service.addRoute(prefix+"/currencies/:base", apiVersion.ListAllConversions)

	service.addRoute(prefix+"/conversion/:from/:to", apiVersion.SingleConversion)
	service.addRoute(prefix+"/convert/:from/:to/:amount", apiVersion.ConvertCurrency)
}

func (service *ApiService) addRoute(path string, handler ControllerHandler) {
	service.server.addHandler("GET", path, func(w http.ResponseWriter, r *http.Request, params denco.Params) {
		Middleware.CORS(service.handleContextRequest(handler))(w, r, params)
	})
}

func (service *ApiService) handleContextRequest(handler ControllerHandler) denco.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request, params denco.Params) {
		ctx := NewRequestContext(r, w, params, service)
		handler(ctx.Get())
	}
}

func (service *ApiService) Listen() {
	fmt.Println("Running on http://127.0.0.1:" + os.Getenv("PORT"))

	http.ListenAndServe(":"+os.Getenv("PORT"), service.server.build())
}
