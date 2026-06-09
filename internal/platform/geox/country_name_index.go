package geox

import (
	"sync"

	"github.com/biter777/countries"
)

var countryNameIndex = struct {
	sync.Once
	values []countryName
}{}

type countryName struct {
	name   string
	alpha2 string
}

func indexedCountryNames() []countryName {
	countryNameIndex.Do(func() {
		countryNameIndex.values = make([]countryName, 0, countries.Total())
		for _, country := range countries.All() {
			if !country.IsValid() {
				continue
			}
			countryNameIndex.values = append(countryNameIndex.values, countryName{
				name:   normalizeCountryNameText(country.String()),
				alpha2: country.Alpha2(),
			})
		}
	})
	return countryNameIndex.values
}
