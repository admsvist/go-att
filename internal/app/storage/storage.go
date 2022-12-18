package storage

import (
	"encoding/csv"
	"github.com/admsvist/go-att/entity"
	"github.com/gocarina/gocsv"
	"log"
	"os"
)

const fileName = "cities.csv"

type Storage struct {
	A      int
	Cities []*entity.City
}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Read() {
	citiesFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer citiesFile.Close()

	r := csv.NewReader(citiesFile)

	if err := gocsv.UnmarshalCSVWithoutHeaders(r, &s.Cities); err != nil { // Load cities from file
		log.Fatalln(err)
	}
}

func (s Storage) Write() {
	citiesFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		log.Fatalln(err)
	}
	defer citiesFile.Close()

	if _, err := citiesFile.Seek(0, 0); err != nil {
		log.Fatalln(err)
	}

	r := csv.NewWriter(citiesFile)

	if err := gocsv.MarshalCSVWithoutHeaders(&s.Cities, r); err != nil { // Load cities from file
		log.Fatalln(err)
	}
}
