package cart

import (
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/shaxbee/butler/internal/session"
)

type Routes struct {
	logger *slog.Logger
}

func (r Routes) Register(mux *http.ServeMux) {
	mux.HandleFunc("GET /cart", r.GetCart)
	mux.HandleFunc("POST /cart", r.AddToCart)
}

func (r Routes) GetCart(rw http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	token := session.GetToken(req)
	if token == nil {
		r.logger.LogAttrs(ctx, slog.LevelError, "session token not found")
		http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}

func (r Routes) AddToCart(rw http.ResponseWriter, req *http.Request) {
}

func (r Routes) handler(c templ.Component) *templ.ComponentHandler {
	return templ.Handler(c, templ.WithErrorHandler(func(req *http.Request, err error) http.Handler {
		r.logger.LogAttrs(req.Context(), slog.LevelError, "render", slog.String("error", err.Error()))

		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		})
	}))
}
