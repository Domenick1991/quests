package handlers

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"quests/internal"
	"quests/services"
	mock_services "quests/services/mocks"
	"testing"
)

func TestHandler_CreateUser(t *testing.T) {
	//загружаем конфиг
	type mockBehavior func(s *mock_services.MockAuthorization, user internal.User)

	type APITestCase struct {
		name         string
		body         string
		user         internal.User
		mockBehavior mockBehavior
		waitStatus   int
		waitResponse string
	}

	tests := []APITestCase{
		{
			name: "Создание пользователя",
			body: `{"username": "test99", "password": "123", "userIsAdmin": false}`,
			user: internal.User{
				Username: "test99",
				Password: "123",
				Isadmin:  false,
			},
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {
				s.EXPECT().CreateUser(user).Return(2, nil)
			},
			waitStatus:   http.StatusOK,
			waitResponse: `{"id:":2}`,
		},
		{
			name:         "Создание пользователя. Без логина 1",
			body:         `{"username": "", "password": "123", "userIsAdmin": false}`,
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {},
			waitStatus:   http.StatusInternalServerError,
			waitResponse: `{"error":"Не указан логин"}`,
		},
		{
			name:         "Создание пользователя. Без логина 2",
			body:         `{"password": "123", "userIsAdmin": false}`,
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {},
			waitStatus:   http.StatusInternalServerError,
			waitResponse: `{"error":"Не указан логин"}`,
		},
		{
			name:         "Создание пользователя. Без пароля 1",
			body:         `{"username": "test1", "password": "", "userIsAdmin": false}`,
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {},
			waitStatus:   http.StatusInternalServerError,
			waitResponse: `{"error":"Не указан пароль"}`,
		},
		{
			name:         "Создание пользователя. Без пароля 2",
			body:         `{"username": "test1","userIsAdmin": false}`,
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {},
			waitStatus:   http.StatusInternalServerError,
			waitResponse: `{"error":"Не указан пароль"}`,
		},
		{
			name: "Создание пользователя. Пустой запрос",
			body: ``,
			user: internal.User{
				Username: "test99",
				Password: "123",
				Isadmin:  false,
			},
			mockBehavior: func(s *mock_services.MockAuthorization, user internal.User) {},
			waitStatus:   http.StatusBadRequest,
			waitResponse: `{"error":"Неверный входной Json"}`,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			auth := mock_services.NewMockAuthorization(c)
			tc.mockBehavior(auth, tc.user)

			s := &services.Service{auth, nil, nil}
			handler := NewHandler(s)
			gin.SetMode(gin.ReleaseMode)
			router := gin.New()

			router.POST("/CreateUser", handler.CreateUser)

			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/CreateUser", bytes.NewBufferString(tc.body))

			router.ServeHTTP(w, r)

			assert.Equal(t, tc.waitStatus, w.Code)
			assert.Equal(t, tc.waitResponse, w.Body.String())

		})
	}
}
