package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

// AdminAuth Авторизация администратора
func (h *Handler) AdminAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if ok {
		user, err := h.services.GetUser(username, password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error:": err.Error()})
			return
		} else {
			if user.Isadmin {
				c.Set("userId", user.Id)
				return
			}
		}
	}
}

// UserAuth Авторизация любого пользователя
func (h *Handler) UserAuth(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if ok {
		user, err := h.services.GetUser(username, password)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error:": err.Error()})
			return
		} else {
			c.Set("userId", user.Id)
			return
		}
	}
}

func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get("userId")
	if !ok {
		return 0, errors.New("userId не указан")
	}

	idInt, ok := id.(int)
	if !ok {
		return 0, errors.New("userId должно быть целым числом")
	}

	return idInt, nil
}
