package root

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/a-h/templ"
	"github.com/bojanz/currency"

	"github.com/shaxbee/butler/assets"
	"github.com/shaxbee/butler/product"
	"github.com/shaxbee/butler/templates/pages"
)

type Routes struct {
	logger   *slog.Logger
	products *product.Service
}

func NewRoutes(logger *slog.Logger, products *product.Service) Routes {
	return Routes{
		logger: logger.With(slog.String("component", "root")),
	}
}

func (r Routes) Register(mux *http.ServeMux) {
	mux.Handle("GET /", r.Home())
	mux.Handle("GET /assets/{asset...}", r.Assets())
}

func (r Routes) Home() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		products, err := r.products.List(req.Context(), "THB")
		r.handler(pages.HomePage(products, err))
	})
}

func (r Routes) Assets() http.Handler {
	handler := http.FileServer(http.FS(assets.FS))
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		req.URL.Path = path.Join("/dist", req.PathValue("asset"))
		handler.ServeHTTP(rw, req)
	})
}

func (r Routes) handler(c templ.Component) *templ.ComponentHandler {
	return templ.Handler(c, templ.WithErrorHandler(func(req *http.Request, err error) http.Handler {
		r.logger.LogAttrs(req.Context(), slog.LevelError, "render", slog.String("error", err.Error()))

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		})
	}))
}

func thb(n string) currency.Amount {
	amount, err := currency.NewAmount(n, "THB")
	if err != nil {
		panic(err)
	}

	return amount
}
