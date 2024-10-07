package server

import (
	"music-library/internal/server/query"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.BaseHandler)

	r.GET("/songs", s.GetSongsHandler)

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
