package entity

import (
	"net/http"
)

type City struct {
	Id         int    `json:"id"`         // (уникальный номер)
	Name       string `json:"name"`       // (название города)
	Region     string `json:"region"`     // (регион)
	District   string `json:"district"`   // (округ)
	Population int    `json:"population"` // (численность населения)
	Foundation int    `json:"foundation"` // (год основания)
}

func (c *City) Bind(r *http.Request) error {
	return nil
}

func (c *City) UpdatePopulation(population int) {
	c.Population = population
}
