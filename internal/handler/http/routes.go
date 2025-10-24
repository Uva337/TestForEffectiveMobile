package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware" 
	httpSwagger "github.com/swaggo/http-swagger"
)


func (h *Handler) RegisterRoutes() http.Handler {
	r := chi.NewRouter()


	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))

	r.Route("/subscriptions", func(r chi.Router) {
		r.Post("/", h.CreateSubscription)
		r.Get("/summary", h.GetSummary)

		r.Route("/{id}", func(r chi.Router) {
			r.Get("/", h.GetSubscriptionByID)
			r.Put("/", h.UpdateSubscription)
			r.Delete("/", h.DeleteSubscription)
		})
	})

	return r
}

