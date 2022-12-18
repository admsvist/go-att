package repository

import (
	"errors"
	"github.com/admsvist/go-att/entity"
	"github.com/admsvist/go-att/internal/app/storage"
)

type Repository struct {
	s *storage.Storage
}

func New(s *storage.Storage) *Repository {
	return &Repository{
		s: s,
	}
}

func (r Repository) GetAll() []*entity.City {
	return r.s.Cities
}

func (r Repository) GetById(cityId int) *entity.City {
	for _, city := range r.s.Cities {
		if city.Id == cityId {
			return city
		}
	}
	return nil
}

func (r Repository) DeleteById(cityId int) {
	for i, city := range r.s.Cities {
		if city.Id == cityId {
			r.s.Cities = append(r.s.Cities[0:i], r.s.Cities[i+1:]...)
			break
		}
	}
}

func (r Repository) Add(city *entity.City) error {
	if r.GetById(city.Id) != nil {
		return errors.New("city with this id already exists")
	}

	r.s.Cities = append(r.s.Cities, city)

	return nil
}
