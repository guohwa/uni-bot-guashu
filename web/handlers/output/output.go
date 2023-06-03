package output

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const VERSION = "v1.0"

func New(ctx *gin.Context) *Output {
	return &Output{
		ctx:     ctx,
		version: VERSION,
	}
}

type Output struct {
	ctx     *gin.Context
	version string
}

func (o *Output) Message(msg interface{}) {
	o.ctx.JSON(http.StatusOK, map[string]interface{}{
		"output":  msg,
		"version": o.version,
	})
}

func (o *Output) Error(err interface{}) {
	o.ctx.JSON(http.StatusOK, map[string]interface{}{
		"error":   err,
		"version": o.version,
	})
}
