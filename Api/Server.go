package Api

import (
	"log"
	"net/http"
	"time"

	"github.com/naoina/denco"
)

type RequestHandler struct {
	handler  http.Handler
	router   *denco.Mux
	handlers []denco.Handler
	Counts   int64
}

func NewServer() *RequestHandler {
	return &RequestHandler{
		router:   denco.NewMux(),
		handlers: []denco.Handler{},
	}
}

func (c *RequestHandler) build() http.Handler {
	handler, err := c.router.Build(c.handlers)
	if err != nil {
		log.Fatal(err)
	}

	c.handler = handler

	return c
}

func (c *RequestHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	start := time.Now()
	defer func() {
		c.Counts++
		diff := time.Now().Sub(start)
		tookMicroS := time.Duration(diff.Microseconds()).String()
		tookMs := time.Duration(diff.Milliseconds()).String()
		log.Printf("Handled request %s in %v (%v) \nTotal Requests: %d \n", request.URL.Path, tookMicroS, tookMs, c.Counts)
	}()

	c.handler.ServeHTTP(writer, request)
}

func (c *RequestHandler) addHandler(method string, path string, handler denco.HandlerFunc) {
	c.handlers = append(c.handlers, c.router.Handler(method, path, handler))
}
