package api

import "github.com/gin-gonic/gin"

type Error struct {
	Message string `json:"message"`
}

type BaseResponse struct {
	Errors []Error     `json:"errors"`
	Data   interface{} `json:"data"`
}

func reject(c *gin.Context, code int, errs ...error) {
	errList := make([]Error, len(errs))
	for i, err := range errs {
		errList[i] = Error{Message: err.Error()}
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