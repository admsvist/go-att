package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"net/http"
	"time"
)

type Handler interface {
	GetAll(http.ResponseWriter, *http.Request)
	GetById(http.ResponseWriter, *http.Request)
	DeleteById(http.ResponseWriter, *http.Request)
	UpdatePopulationById(http.ResponseWriter, *http.Request)
	Create(http.ResponseWriter, *http.Request)
}

type Router struct {
	r *chi.Mux
}

func New(h Handler) *Router {
	r := chi.NewRouter()

	r.Use(
		middleware.RedirectSlashes,
		middleware.Timeout(30*time.Second),
		render.SetContentType(render.ContentTypeJSON),
	)

	r.Route("/cities", func(r chi.Router) {
		r.Get("/", h.GetAll)  // получить все города
		r.Post("/", h.Create) // добавить новый город
		r.Route("/{cityId}", func(r chi.Router) {
			r.Get("/", h.GetById)                          // получить город по ид
			r.Delete("/", h.DeleteById)                    // удалить город по ид
			r.Patch("/population", h.UpdatePopulationById) // обновить население по ид
		})
	})

	return &Router{
		r: r,
	}
}

func (r Router) GetRouter() *chi.Mux {
	return r.r
}
