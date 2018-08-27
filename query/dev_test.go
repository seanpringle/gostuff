package query

import (
	"log"
	"testing"
)

const (
	Name Field = Id + 1 + iota
	Country
)

func Test1(test *testing.T) {

	countries := NewTable()
	au := countries.Insert(Tuple{
		(Name): "Australia",
	})
	uk := countries.Insert(Tuple{
		(Name): "England",
	})
	countries.Commit()

	cities := NewTable()
	cities.Insert(Tuple{
		(Name):    "Canberra",
		(Country): au,
	})
	cities.Insert(Tuple{
		(Name):    "Sydney",
		(Country): au,
	})
	cities.Insert(Tuple{
		(Name):    "London",
		(Country): uk,
	})
	cities.Insert(Tuple{
		(Name):    "Edinburgh",
		(Country): uk,
	})
	cities.Commit()

	Save("countries", countries)
	Save("cities", cities)

	ccs := Select(countries).In(Name, "England", "Australia").List(Id)

	for tuple := range Select(cities, Name).In(Country, ccs...).Run() {
		log.Println(tuple)
	}

	countries, err1 := Load("countries")
	log.Println(countries, err1)

	cities, err2 := Load("cities")
	log.Println(cities, err2)

	for tuple := range Select(cities).Order(Name, ASC).Run() {
		log.Println(tuple)
	}

	log.Println("group")

	ids := Select(cities).Min(Id).Group(Country).List(Id)
	for tuple := range Select(cities).In(Id, ids...).Run() {
		log.Println(tuple)
	}
}
