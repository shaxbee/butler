package root

import (
	"log/slog"
	"net/http"
	"path"

	"github.com/a-h/templ"
	"github.com/bojanz/currency"

	"github.com/shaxbee/butler/assets"
	"github.com/shaxbee/butler/internal/product"
	"github.com/shaxbee/butler/templates/pages"
)

type Routes struct {
	logger *slog.Logger
}

func NewRoutes(logger *slog.Logger) Routes {
	return Routes{
		logger: logger.With(slog.String("component", "root")),
	}
}

func (r Routes) Register(mux *http.ServeMux) {
	mux.Handle("GET /", r.Home())
	mux.Handle("GET /assets/{asset...}", r.Assets())
}

func (r Routes) Home() http.Handler {
	products := []product.Product{
		{
			Name:        "Product 1",
			Category:    "Lorem",
			Price:       thb("500"),
			Description: "lorem, ipsum, dolor",
		},
		{
			Name:            "Product 2",
			Category:        "Lorem",
			Price:           thb("800"),
			DiscountedPrice: thb("600"),
			Description:     "lorem, ipsum, dolor",
		},
		{
			Name:        "Product 3",
			Category:    "Lorem",
			Price:       thb("700"),
			Description: "elit, sed, do",
		},
		{
			Name:            "Product 4",
			Category:        "Ipsum",
			Price:           thb("700"),
			DiscountedPrice: thb("450"),
			Description:     "amet, consectetur, adipiscing",
		},
		{
			Name:        "Product 5",
			Category:    "Ipsum",
			Price:       thb("500"),
			Description: "do, eiusmod, tempor",
		},
	}
	return r.handler(pages.HomePage(products))
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
		r.logger.Error("render", slog.String("error", err.Error()))
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
