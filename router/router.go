package router

import (
	"github.com/gin-gonic/gin"
	"wq_submitter/internal/auth"

	"wq_submitter/api"
)

func SetRouter() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	//对所有请求使用新的 API 密钥认证中间件
	r.Use(auth.APIKeyAuthMiddleware())

	wqSubmitter := r.Group("/wq_submitter")
	wqSubmitter.GET("/hello", api.Hello)
	ideaGroup := wqSubmitter.Group("/idea")
	ideaGroup.GET("/all", api.GetAllIdea)
	ideaGroup.GET("/unfinish", api.GetUnfinishedIdea)
	ideaGroup.GET("/run", api.GetRunningIdea)
	ideaGroup.POST("/concurrency", api.UpdateIdeaConcurrency)
	ideaGroup.POST("/delete", api.DeleteIdea)

	alphaGroup := wqSubmitter.Group("/alpha")

	alphaGroup.POST("/upload", api.UploadAlphaListWithIdea)
	alphaGroup.GET("/list", api.GetAlphaListByIdea)
	return r
}
