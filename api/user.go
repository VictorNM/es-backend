package api

import (
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/auth"
	"github.com/victornm/es-backend/store/memory"
	"github.com/victornm/es-backend/user"
	"log"
	"net/http"
)

func (s *Server) createGetProfileHandler() func(c *gin.Context) {
	userQuery := s.createUserGetProfileQuery()

	return func(c *gin.Context) {
		userAuth, ok := c.Get("user")
		if !ok {
			log.Panic("something wrong")
		}

		u, err := userQuery.GetProfile(userAuth.(*auth.UserAuth).UserID)
		if err != nil {
			reject(c, http.StatusNotFound, err)
		}

		response(c, http.StatusOK, u)
	}
}

func (s *Server) createUserGetProfileQuery() user.GetProfileQuery {
	return user.NewQueryService(memory.NewUserStore())
}