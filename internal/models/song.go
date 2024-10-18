package models

type Song struct {
	Id          int    `json:"id"`
	Group       string `json:"group"`
	Song        string `json:"song"`
	ReleaseDate Date   `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type NewSong struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
