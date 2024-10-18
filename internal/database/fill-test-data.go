package database

import (
	"log"
	"music-library/internal/models"
)

func FillTestData(s *service) {
	songs := []models.Song{
		{
			Group:       "The Rolling Stones",
			Song:        "Satisfaction",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Satisfaction by The Rolling Stones",
			Link:        "https://en.wikipedia.org/wiki/Satisfaction_(song)",
		},
		{
			Group:       "The Rolling Stones",
			Song:        "Paint it black",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Paint it black by The Rolling Stones",
			Link:        "https://en.wikipedia.org/wiki/Paint_it_black",
		},
		{
			Group:       "The Rolling Stones",
			Song:        "Shake it off",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Shake it off by The Rolling Stones",
			Link:        "https://en.wikipedia.org/wiki/Shake_it_off",
		},
		{
			Group:       "Adele",
			Song:        "Rolling in the deep",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Rolling in the deep by The Adele",
			Link:        "https://en.wikipedia.org/wiki/Rolling_in_the_deep",
		},
		{
			Group:       "Adele",
			Song:        "Someone like you",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Someone like you by The Adele",
			Link:        "https://en.wikipedia.org/wiki/Something_Like_You",
		},
		{
			Group:       "Muse",
			Song:        "Darkshines",
			ReleaseDate: models.NewDateFromString("1978-08-10"),
			Text:        "Darkshines by Muse",
			Link:        "https://en.wikipedia.org/wiki/Darkshines",
		},
	}
	for _, song := range songs {
		err := s.AddNewSong(song)
		if err != nil {
			log.Fatalf("Can't add new song : %v\n", err)
		}
	}
}
