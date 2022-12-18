package handler

import (
	"encoding/json"
	"errors"
	"github.com/admsvist/go-att/entity"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"net/http"
	"strconv"
	"strings"
)

type Repository interface {
	GetAll() []*entity.City
	GetById(int) *entity.City
	DeleteById(int)
	Add(*entity.City) error
}

type Handler struct {
	r Repository
}

func New(r Repository) *Handler {
	return &Handler{
		r: r,
	}
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func StatusError(status int, err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: status,
		StatusText:     err.Error(),
		ErrorText:      err.Error(),
	}
}

func (h Handler) GetById(w http.ResponseWriter, r *http.Request) {
	cityIdString := chi.URLParam(r, "cityId")

	cityId, err := strconv.Atoi(cityIdString)
	if err != nil {
		render.Render(w, r, StatusError(http.StatusInternalServerError, err))
		return
	}

	city := h.r.GetById(cityId)
	if city == nil {
		render.Render(w, r, StatusError(http.StatusNotFound, errors.New("404 not found")))
		return
	}

	render.JSON(w, r, city)
	render.Status(r, http.StatusOK)
}

func (h Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	cities := h.r.GetAll()

	filterByRegion := func(r *http.Request, cities *[]*entity.City) {
		region := r.URL.Query().Get("region")
		if region == "" {
			return
		}

		var filtered []*entity.City
		for _, city := range *cities {
			if city.Region == region {
				filtered = append(filtered, city)
			}
		}
		*cities = filtered
	}

	filterByDistrict := func(r *http.Request, cities *[]*entity.City) {
		district := r.URL.Query().Get("district")
		if district == "" {
			return
		}

		var filtered []*entity.City
		for _, city := range *cities {
			if city.District == district {
				filtered = append(filtered, city)
			}
		}
		*cities = filtered
	}

	filterByPopulationRange := func(r *http.Request, cities *[]*entity.City) error {
		population := r.URL.Query().Get("population")
		if population == "" {
			return nil
		}

		borders := strings.Split(population, "-")
		if len(borders) != 2 {
			return errors.New("population value must be in the format 0-10")
		}

		lower, err := strconv.Atoi(borders[0])
		if err != nil {
			return err
		}

		upper, err := strconv.Atoi(borders[1])
		if err != nil {
			return err
		}

		var filtered []*entity.City
		for _, city := range *cities {
			if city.Population >= lower && city.Population <= upper {
				filtered = append(filtered, city)
			}
		}
		*cities = filtered

		return nil
	}

	filterByFoundationRange := func(r *http.Request, cities *[]*entity.City) error {
		foundation := r.URL.Query().Get("foundation")
		if foundation == "" {
			return nil
		}

		borders := strings.Split(foundation, "-")
		if len(borders) != 2 {
			return errors.New("foundation value must be in the format 0-10")
		}

		lower, err := strconv.Atoi(borders[0])
		if err != nil {
			return err
		}

		upper, err := strconv.Atoi(borders[1])
		if err != nil {
			return err
		}

		var filtered []*entity.City
		for _, city := range *cities {
			if city.Foundation >= lower && city.Foundation <= upper {
				filtered = append(filtered, city)
			}
		}
		*cities = filtered

		return nil
	}

	filterByRegion(r, &cities)
	filterByDistrict(r, &cities)

	if err := filterByPopulationRange(r, &cities); err != nil {
		render.Render(w, r, StatusError(http.StatusBadRequest, err))
		return
	}

	if err := filterByFoundationRange(r, &cities); err != nil {
		render.Render(w, r, StatusError(http.StatusBadRequest, err))
		return
	}

	render.JSON(w, r, cities)
	render.Status(r, http.StatusOK)
}

func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	city := &entity.City{}
	if err := render.Bind(r, city); err != nil {
		render.Render(w, r, StatusError(http.StatusBadRequest, err))
		return
	}

	if err := h.r.Add(city); err != nil {
		render.Render(w, r, StatusError(http.StatusBadRequest, err))
		return
	}

	render.JSON(w, r, city)
	render.Status(r, http.StatusCreated)
}

func (h Handler) DeleteById(w http.ResponseWriter, r *http.Request) {
	cityIdString := chi.URLParam(r, "cityId")

	cityId, err := strconv.Atoi(cityIdString)
	if err != nil {
		render.Render(w, r, StatusError(http.StatusInternalServerError, err))
		return
	}

	city := h.r.GetById(cityId)
	if city == nil {
		render.Render(w, r, StatusError(http.StatusNotFound, errors.New("404 not found")))
		return
	}

	h.r.DeleteById(cityId)

	render.Status(r, http.StatusNoContent)
}

func (h Handler) UpdatePopulationById(w http.ResponseWriter, r *http.Request) {
	cityIdString := chi.URLParam(r, "cityId")

	cityId, err := strconv.Atoi(cityIdString)
	if err != nil {
		render.Render(w, r, StatusError(http.StatusInternalServerError, err))
		return
	}

	city := h.r.GetById(cityId)
	if city == nil {
		render.Render(w, r, StatusError(http.StatusNotFound, errors.New("404 not found")))
		return
	}

	p := map[string]interface{}{}
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		render.Render(w, r, StatusError(http.StatusInternalServerError, err))
		return
	}

	population, ok := p["population"]
	if !ok {
		render.Render(w, r, StatusError(http.StatusBadRequest, errors.New("missing field population")))
		return
	}

	population, ok = population.(float64)
	if !ok {
		render.Render(w, r, StatusError(http.StatusBadRequest, errors.New("population should be a number")))
		return
	}

	city.UpdatePopulation(int(population.(float64)))

	render.JSON(w, r, city)
	render.Status(r, http.StatusOK)
}
