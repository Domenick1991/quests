package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	_ "quests/docs"
	"quests/services"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	adminApi := router.Group("/Users", h.AdminAuth)
	{
		adminApi.POST("/CreateUser", h.CreateUser)
		adminApi.DELETE("/DeleteUser", h.DeleteUser)
	}

	userApi := router.Group("/Users", h.UserAuth)
	{
		userApi.GET("/GetAllUsers", h.GetAllUser)
	}

	questApi := router.Group("/Quests", h.UserAuth)
	{
		questApi.GET("/GetAllUsers", h.GetAllUser)
		questApi.GET("/GetQuests", h.GetQuests)
		questApi.POST("/CreateUser", h.CreateUser)
		questApi.POST("/CreateQuestSteps", h.CreateQuestSteps)
		questApi.POST("/UpdateQuestSteps", h.UpdateQuestSteps)
		questApi.POST("/CreateQuest", h.CreateQuest)
	}

	historyApi := router.Group("/History", h.UserAuth)
	{
		historyApi.GET("/GetHistory", h.GetHistory)
		historyApi.POST("/CompleteSteps", h.CompleteSteps)
	}

	return router
}
