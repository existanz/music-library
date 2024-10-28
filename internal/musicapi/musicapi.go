package musicapi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"music-library/internal/models"
)

func GetMusicInfo(group_name, song_name string) (models.Song, error) {
	song := models.Song{
		Group: group_name,
		Song:  song_name,
	}
	query := fmt.Sprintf("%s?group=%s&song=%s", os.Getenv("EXTERNAL_API_URL"), url.QueryEscape(group_name), url.QueryEscape(song_name))

	resp, err := http.Get(query)
	if err != nil {
		return song, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil || body == nil {
		return song, err
	}

	if err := json.Unmarshal(body, &song); err != nil {
		return song, err
	}

	if err := resp.Body.Close(); err != nil {
		return song, err
	}

	if resp.StatusCode != http.StatusOK {
		return song, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return song, nil
}
