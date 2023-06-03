package response

import (
	"net/http"
	"regexp"

	"bot/log"

	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

const pattern = "^[a-zA-Z0-9_]+$"

var matcher = regexp.MustCompile(pattern)

func New(ctx *gin.Context) *Response {
	return &Response{
		ctx: ctx,
	}
}

type Context map[string]interface{}

type AjaxResponse struct {
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
	URL  string      `json:"url"`
	Wait int         `json:"wait"`
}

type TableData struct {
	Code  int         `json:"code"`
	Count int64       `json:"count"`
	Items interface{} `json:"items"`
}

type Response struct {
	ctx *gin.Context
}

func (r *Response) Error(msg interface{}) {
	if err, ok := msg.(error); ok {
		log.Error(err)
		r.ctx.JSON(http.StatusOK, AjaxResponse{
			Code: 1,
			Msg:  err.Error(),
			Wait: 2,
		})
		return
	}

	r.ctx.JSON(http.StatusOK, AjaxResponse{
		Code: 1,
		Msg:  msg,
		Wait: 2,
	})
}

func (r *Response) Success(msg interface{}, url string) {
	r.ctx.JSON(http.StatusOK, AjaxResponse{
		Msg:  msg,
		URL:  url,
		Wait: 2,
	})
}

func (r *Response) HTML(name string, data Context) {
	r.ctx.HTML(http.StatusOK, name, r.with(data))
}

func (r *Response) String(data string) {
	r.ctx.String(http.StatusOK, data)
}

func (r *Response) Unauthorized(msg interface{}) {
	r.ctx.JSON(http.StatusUnauthorized, AjaxResponse{
		Msg:  "Please Login",
		URL:  "/account/login",
		Wait: 2,
	})
}

func (r *Response) Table(data TableData) {
	r.ctx.JSON(http.StatusOK, data)
}

func (r *Response) with(data Context) pongo2.Context {
	context := pongo2.Context{}
	for k, v := range r.ctx.Keys {
		if matcher.MatchString(k) {
			context[k] = v
		}
	}

	for k, v := range data {
		if matcher.MatchString(k) {
			context[k] = v
		}
	}

	return context
}
