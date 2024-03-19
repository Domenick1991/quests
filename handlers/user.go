package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"quests/internal"
)

// @Summary Создать пользователя
// @Tags user
// @Description Создает нового пользователя приложения
// @id CreateUser
// @Accept json
// @Procedure json
// @param input body internal.User true "Информация о пользователе"
// @router /Users/CreateUser [post]
// @Security BasicAuth
func (h *Handler) CreateUser(c *gin.Context) {
	var user internal.User

	if err := c.ShouldBindJSON(&user); err == nil {
		if user.Username == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Не указан логин"})
			return
		}
		if user.Password == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Не указан пароль"})
			return
		}
		id, err := h.services.CreateUser(user)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			c.JSON(http.StatusOK, gin.H{"id:": id})
		}
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}

}

// @Summary Показать всех пользователей
// @Tags user
// @Description Показать всех пользователей
// @id GetAllUser
// @Accept json
// @Procedure json
// @router /Users/GetAllUser [GET]
// @Success 200 {object} []internal.User
// @Security BasicAuth
func (h *Handler) GetAllUser(c *gin.Context) {
	c.JSON(http.StatusOK, h.services.GetAllUser())
}

// @Summary Удалить пользователя
// @Tags user
// @Description Удаляет пользователя приложения
// @id DeleteUser
// @Accept json
// @Procedure json
// @param input body internal.DeleteUserStruct true "Идентификатор пользователя"
// @router /Users/DeleteUser [Delete]
// @Security BasicAuth
func (h *Handler) DeleteUser(c *gin.Context) {
	var user internal.DeleteUserStruct
	if err := c.ShouldBindJSON(&user); err == nil {
		err = h.services.DeleteUser(user.Id)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "не удалось удалить пользователя"})
		}
		c.JSON(http.StatusOK, gin.H{"Успешно": "Пользователь удален"})
	} else {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"Неверный входной Json": err.Error()})
		return
	}

}
