package server

import (
	"music-library/internal/musicapi"
	"music-library/internal/server/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.BaseHandler)

	r.GET("/songs", s.GetSongsHandler)

	r.POST("/songs", s.AddNewSongHandler)

	return r
}

func (s *Server) BaseHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "This is the base handler"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) GetSongsHandler(c *gin.Context) {
	data, err := s.db.GetSongs(query.GetOptions(c))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

func (s *Server) AddNewSongHandler(c *gin.Context) {
	var newSong struct {
		Group_name string `json:"group"`
		Song_name  string `json:"song"`
	}

	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	song, err := musicapi.GetMusicInfo(newSong.Group_name, newSong.Song_name)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	err = s.db.AddNewSong(song)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, "Song added")
}
