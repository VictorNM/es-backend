package api

import (
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/store/memory"
	"github.com/victornm/es-backend/user"
	"net/http"
)

// @Summary Get current sign-inned user's profile
// @Description Get profile by user_id in token,
// @Tags user
// @Produce json
// @Success 200 {object} api.BaseResponse{data=user.ProfileDTO} "Get profile successfully"
// @Router /users/profile [get]
func (s *Server) createGetProfileHandler() func(c *gin.Context) {
	userQuery := s.createUserGetProfileQuery()

	return func(c *gin.Context) {
		userAuth := getUser(c)

		u, err := userQuery.GetProfile(userAuth.UserID)
		if err != nil {
			reject(c, http.StatusNotFound, err)
		}

		response(c, http.StatusOK, u)
	}
}

func (s *Server) createUserGetProfileQuery() user.GetProfileQuery {
	return user.NewQueryService(memory.UserStore)
}
