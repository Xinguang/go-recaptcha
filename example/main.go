package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xinguang/go-recaptcha"
	"strings"
)

type JsonBody struct {
	Token             string `json:"token" form:"token"`
	RecaptchaResponse string `json:"g-recaptcha-response" form:"g-recaptcha-response"`
}

func main() {
	router := gin.Default()
	router.Static("/", "./html")

	captcha, err := recaptcha.NewReCAPTCHA()

	router.POST("/signin", func(c *gin.Context) {
		var body JsonBody
		contentType := c.Request.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			c.Bind(&body)
			logrus.Info("form:", body)
		} else {
			c.BindJSON(&body)
			logrus.Info("json:", body)
		}
		token := body.Token
		if len(token) == 0 {
			token = body.RecaptchaResponse
		}
		logrus.Info("token:", token)

		err = captcha.Verify(token)
		if err != nil {
			c.JSON(401, gin.H{
				"error": err,
				"token": token,
				"c":     c.PostForm("token"),
			})
			return
		}
		c.JSON(200, gin.H{
			"token":   token,
			"message": "valid",
		})
	})
	// Listen and serve on 0.0.0.0:8002
	router.Run(":8002")
}
