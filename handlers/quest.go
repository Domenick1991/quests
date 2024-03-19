package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quests/internal"
)

// @Summary Показывает все задания
// @Tags quests
// @Description Показывает все задания в системе
// @id GetQuests
// @Accept json
// @Procedure json
// @router /Quests/GetQuests [GET]
// @Success 200 {object} []internal.Quests
// @Security BasicAuth
func (h *Handler) GetQuests(c *gin.Context) {
	quests, err := h.services.GetQuests()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Ошибка при получении заданий": err.Error()})
		return
	}
	c.JSON(http.StatusOK, quests)
}

// @Summary Добавить задание
// @Tags quests
// @Description Создает новое задание
// @id CreateQuest
// @Accept json
// @Procedure json
// @router /Quests/CreateQuest [POST]
// @param input body internal.NewQuest true "информация о задании"
// @Success 200
// @Security BasicAuth
func (h *Handler) CreateQuest(c *gin.Context) {

	var quest internal.NewQuest

	if err := c.ShouldBindJSON(&quest); err == nil {

		errlist := h.services.CreateQuest(quest)
		if errlist != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errlist)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"Успешно:": "Задания успешно созданы"})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}
}

// @Summary Добавляет информацию о новых шагах к заданию
// @Tags quests
// @Description Добавляет информацию о новых шагах к заданию
// @id CreateQuestSteps
// @Accept json
// @Procedure json
// @router /Quests/CreateQuestSteps [POST]
// @param input body internal.NewQuestSteps true "информация о шагах задания"
// @Success 200
// @Security BasicAuth
func (h *Handler) CreateQuestSteps(c *gin.Context) {
	var questSteps internal.NewQuestSteps
	if err := c.ShouldBindJSON(&questSteps); err == nil {

		errlist := h.services.CreateQuestSteps(questSteps)
		if errlist != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errlist)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"Успешно:": "шаги к заданию успешно созданы"})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}

}

// @Tags quests
// @Description Обновляет информацию о шагах заданий
// @id UpdateQuestSteps
// @Accept json
// @Procedure json
// @router /Quests/UpdateQuestSteps [POST]
// @param input body internal.UpdateQuestSteps true "обновленная информация о шагах задания"
// @Success 200
// @Security BasicAuth
func (h *Handler) UpdateQuestSteps(c *gin.Context) {

	var updateQuestSteps internal.UpdateQuestSteps

	if err := c.ShouldBindJSON(&updateQuestSteps); err == nil {

		errlist := h.services.UpdateQuestSteps(updateQuestSteps)
		if errlist != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, errlist)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"Успешно:": "шаги к заданию успешно обновлены"})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}

}
