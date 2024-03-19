package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quests/internal"
)

// @Summary История по пользователю
// @Tags history
// @Description Показывает историю выполнения заданий для пользователя
// @id GetHistory
// @Accept json
// @Procedure json
// @router /History/GetHistory [GET]
// @Success 200 {object} internal.UserBonus
// @Security BasicAuth
func (h *Handler) GetHistory(c *gin.Context) {
	userId, err := getUserId(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	} else {
		userbous, err := h.services.GetHistory(userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Ошибка при получении данных": err.Error()})
			return
		}
		c.JSON(http.StatusOK, userbous)
	}
}

// @Summary установить пометки у выполнении для шагов
// @Tags history
// @Description Устанавливает признак выполнения шага для пользователя
// @id CompleteSteps
// @Accept json
// @Procedure json
// @router /History/CompleteSteps [POST]
// @param input body internal.NewCompleteSteps true "обновленная информация о шагах задания"
// @Success 200 {object} internal.NewCompleteSteps
// @Security BasicAuth
func (h *Handler) CompleteSteps(c *gin.Context) {
	var questSteps internal.NewCompleteSteps
	if err := c.ShouldBindJSON(&questSteps); err == nil {
		errlist := h.services.CompleteSteps(questSteps)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, errlist)
			return
		}
		c.JSON(http.StatusOK, questSteps)
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}

}
