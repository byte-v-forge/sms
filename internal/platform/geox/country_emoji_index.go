package geox

import (
	"sync"

	"github.com/biter777/countries"
)

var countryEmojiIndex = struct {
	sync.Once
	values []countryEmoji
}{}

type countryEmoji struct {
	emoji  string
	alpha2 string
}

func countryEmojis() []countryEmoji {
	countryEmojiIndex.Do(func() {
		countryEmojiIndex.values = make([]countryEmoji, 0, countries.Total())
		for _, country := range countries.All() {
			if !country.IsValid() {
				continue
			}
			countryEmojiIndex.values = append(countryEmojiIndex.values, countryEmoji{
				emoji:  country.Emoji(),
				alpha2: country.Alpha2(),
			})
		}
	})
	return countryEmojiIndex.values
}
