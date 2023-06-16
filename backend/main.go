package main

import (
	"backend/common"
	_ "backend/common"
	conf "backend/config"
	"backend/service/user"
	"backend/service/word"
	"backend/utils"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

var (
	z = common.GetLogger()
)

func main() {
	router := gin.Default()

	v1 := router.Group("/login-api")
	v2 := router.Group("/basic-api")
	v3 := router.Group("/batch-api")

	v1.Use()
	{
		v1.POST("/login", user.Login)
		v1.POST("/register", user.Register)
		v1.POST("/resetPassword", user.ResetPassword)
		v1.GET("/generateCode", user.GenerateCaptcha)
		v1.POST("/verifyCode", user.VerifyCode)
	}

	v2.Use(utils.AuthUser())
	{

		v2.POST("/setGrade", user.SetGrade)
		v2.GET("/getUserInfo", user.GetUserInfo)

		v2.GET("/word/getList1", word.GetWordList1)
		v2.POST("/word/postList1", word.JudgeUserWordLevel)
		v2.POST("/word/getVocabulary", word.CalculateVocabulary)
	}

	v3.POST("/word/batchProcess", word.BatchProcess)

	server := &http.Server{
		Addr:           ":" + conf.GlobalConfig.WebSever.HttpsListenPort,
		Handler:        router,
		ReadTimeout:    3600 * time.Second,
		WriteTimeout:   3600 * time.Second,
		MaxHeaderBytes: 32 << 20,
	}
	z.Info("listening" + conf.GlobalConfig.WebSever.HttpsListenPort)
	err := server.ListenAndServe()
	if err != nil {
		log.Println("服务启动失败!", err)
	}
}
