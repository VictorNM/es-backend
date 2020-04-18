package api

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/victornm/es-backend/pkg/auth"
	"log"
)

type Error struct {
	Message string `json:"message"`
	Detail  string `json:"detail"`
}

type BaseResponse struct {
	Errors []Error     `json:"errors"`
	Data   interface{} `json:"data"`
}

func reject(c *gin.Context, code int, errs ...error) {
	errList := make([]Error, len(errs))
	for i, err := range errs {
		if err == nil {
			continue
		}

		if unwrap := errors.Unwrap(err); unwrap != nil {
			errList[i] = Error{
				Message: unwrap.Error(),
				Detail:  err.Error(),
			}
			continue
		}

		errList[i] = Error{
			Message: err.Error(),
			Detail:  err.Error(),
		}
	}

	c.JSON(code, &BaseResponse{
		Errors: errList,
		Data:   nil,
	})
}

func response(c *gin.Context, code int, data interface{}) {
	c.JSON(code, &BaseResponse{
		Errors: nil,
		Data:   data,
	})
}

func getUser(c *gin.Context) *auth.UserAuthDTO {
	userAuth, ok := c.Get("user")
	if !ok {
		log.Panic("key 'user' should be present in context")
	}

	return userAuth.(*auth.UserAuthDTO)
}
