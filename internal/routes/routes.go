package routes

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/julienschmidt/httprouter"
	"github.com/shaxbee/butler/templates/pages"
)

type Routes struct{}

func (r Routes) Handler() http.Handler {
	router := httprouter.New()
	router.Handler(http.MethodGet, "/", r.Home())

	return router
}

func (r Routes) Home() http.Handler {
	return templ.Handler(pages.HomePage())
}
