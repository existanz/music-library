package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"music-library/internal/customErrors"
	"music-library/internal/models"
	"music-library/internal/musicapi"
	"music-library/internal/server/query"

	_ "music-library/docs"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	swagger "github.com/swaggo/gin-swagger"
)

const verseDelimiter = "\n\n"

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/docs/*any", swagger.WrapHandler(swaggerfiles.Handler))

	r.GET("/songs", s.GetSongsHandler)

	r.GET("/songs/:id", s.GetSongByIdHandler)

	r.GET("/songs/:id/:verse", s.GetSongTextByVerseHandler)

	r.POST("/songs", s.AddNewSongHandler)

	r.PUT("/songs/:id", s.UpdateSongHandler)

	r.DELETE("/songs/:id", s.DeleteSongHandler)

	return r
}

// GetSongByIdHandler
//
// @Summary		Get song by id
// @Description	Get song by id
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Song ID"
// @Success		200	{object}	models.Song
// @Failure		404	{string}	string	"Song not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/songs/{id} [get]
func (s *Server) GetSongByIdHandler(c *gin.Context) {
	data, err := s.db.GetSongById(c.Param("id"))
	slog.Info("GetSongByIdHandler", "group", data.Group, "song", data.Song)
	if err != nil {
		if err == customErrors.ErrNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// GetSongTextByVerseHandler
//
// @Summary		Get song text by verse
// @Description	Get song text by verse
// @Accept			json
// @Produce		json
// @Param			id		path		int	true	"Song ID"
// @Param			verse	path		int	true	"Verse number"
// @Success		200		{string}	string
// @Failure		404		{string}	string	"Song not found"
// @Failure		500		{string}	string	"Internal server error"
// @Router			/songs/{id}/{verse} [get]
func (s *Server) GetSongTextByVerseHandler(c *gin.Context) {
	data, err := s.db.GetSongById(c.Param("id"))
	if err != nil {
		if err == customErrors.ErrNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	text := strings.Split(data.Text, verseDelimiter)

	slog.Info(fmt.Sprintf("Lirics of the song %s has %d verses", data.Song, len(text)))

	verse, err := strconv.Atoi(c.Param("verse"))
	if err != nil || verse > len(text) {
		verse = 1
	}
	c.String(http.StatusOK, text[verse-1])
}

// GetSongsHandler
//
// @Summary		Get all songs
// @Description	Get all songs
// @Accept			json
// @Produce		json
// @Success		200	{object}	[]models.Song
// @Failure		500	{string}	string	"Internal server error"
// @Router			/songs [get]
func (s *Server) GetSongsHandler(c *gin.Context) {
	data, err := s.db.GetSongs(query.GetOptions(c))
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, data)
}

// AddNewSongHandler
//
//	@Summary		Add new song
//	@Description	Add new song
//	@ID				id
//	@Accept			json
//	@Produce		json
//	@Param			song	body		models.NewSong	true	"Song"
//	@Success		200		{string}	string
//	@Failure		400		{string}	string	"Bad request"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/songs [post]
func (s *Server) AddNewSongHandler(c *gin.Context) {
	var newSong models.NewSong

	if err := c.ShouldBindJSON(&newSong); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	song, err := musicapi.GetMusicInfo(newSong.Group, newSong.Song)
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

// UpdateSongHandler
//
// @Summary		Update song
// @Description	Update song
// @Accept			json
// @Produce		json
// @Param			song	body		models.Song	true	"Song"
// @Success		200		{string}	string
// @Failure		400		{string}	string	"Bad request"
// @Failure		404		{string}	string	"Song not found"
// @Failure		500		{string}	string	"Internal server error"
// @Router			/songs/{id} [put]
func (s *Server) UpdateSongHandler(c *gin.Context) {
	songID := c.Param("id")

	song, err := s.db.GetSongById(songID)
	if err != nil {
		if err == customErrors.ErrNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	if err := c.ShouldBindJSON(&song); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	if song.Id != 0 && songID != strconv.Itoa(song.Id) {
		c.String(http.StatusBadRequest, "Wrong id")
		return
	}
	err = s.db.UpdateSongById(songID, song)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.String(http.StatusOK, fmt.Sprintf("Song id:%s updated", songID))
}

// DeleteSongHandler
//
// @Summary		Delete song
// @Description	Delete song
// @Accept			json
// @Produce		json
// @Param			id	path		int	true	"Song ID"
// @Success		200	{string}	string
// @Failure		404	{string}	string	"Not found"
// @Failure		500	{string}	string	"Internal server error"
// @Router			/songs/{id} [delete]
func (s *Server) DeleteSongHandler(c *gin.Context) {
	songID := c.Param("id")
	err := s.db.DeleteSongById(songID)
	if err != nil {
		if err == customErrors.ErrNotFound {
			c.String(http.StatusNotFound, err.Error())
			return
		}
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, fmt.Sprintf("Song id:%s deleted", songID))
}
