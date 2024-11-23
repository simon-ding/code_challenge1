package server

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"testing"
)

func TestHttpHandler(t *testing.T) {
	HttpHandler(func(context *gin.Context) (interface{}, error) {
		return nil, nil
	})
	HttpHandler(func(context *gin.Context) (interface{}, error) {
		return nil, errors.Errorf("test")
	})

}
