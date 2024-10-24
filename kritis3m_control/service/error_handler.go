package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

//create function ErrorHandler with argument zerologger and return gin.HandlerFunc

type ErrorJson struct {
	status     string `json:"status"`
	error_type string `json:"type"`
}

func ErrorHandler_v1(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		for _, err := range c.Errors {
			switch err.Err {

			}
			// etc...
		}

		c.JSON(http.StatusInternalServerError, "")
	}
}
