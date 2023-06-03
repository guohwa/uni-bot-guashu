package handlers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/uncle-gua/base64Captcha"
)

var (
	driver         = base64Captcha.DefaultDriverDigit
	store          = base64Captcha.DefaultMemStore
	captchaHandler = &captcha{}
)

type captcha struct {
}

func (handler *captcha) Handle(router *gin.Engine) {
	router.GET("/captcha", func(ctx *gin.Context) {
		key, content, answer := driver.GenerateIdQuestionAnswer()
		item, err := driver.DrawCaptcha(content)
		if err != nil {
			panic(err)
		}

		store.Set(key, answer)

		session := sessions.Default(ctx)
		session.Set("captcha-key", key)
		session.Save()

		ctx.Header("Content-Type", "image/png")
		item.WriteTo(ctx.Writer)
	})
}

func (handler *captcha) Verify(ctx *gin.Context, code string) bool {
	session := sessions.Default(ctx)
	key := session.Get("captcha-key")
	if s, ok := key.(string); ok {
		return store.Verify(s, code, true)
	}

	return false
}
