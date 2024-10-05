package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.GET("/", s.BaseHandler)

	return r
}

func (s *Server) BaseHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "This is the base handler"

	c.JSON(http.StatusOK, resp)
}
