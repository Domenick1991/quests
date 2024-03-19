package repository

import (
	"errors"
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"quests/internal"
)

type AuthRepo struct {
	db *dbx.DB
}

func NewAuthRepo(db *dbx.DB) *AuthRepo {
	return &AuthRepo{db: db}
}

func (auth *AuthRepo) CreateUser(user internal.User) (int, error) {
	//хешируем пароль
	user.Password = EncodePassword(user.Password)

	//Проверяем существует ли такой пользователь по имени
	var oldUser internal.User
	err := auth.db.Select().From(user.TableName()).Where(dbx.HashExp{"username": user.Username}).One(&oldUser)
	if err != nil { //ошибка в случае dbx возвращается если значения не нашлись
		err1 := auth.db.Model(&user).Insert()
		if err1 != nil {
			return 0, errors.New("не удалось добавить пользователя")
		} else {
			return user.Id, nil
		}
	} else {
		return 0, fmt.Errorf("пользователь уже существует. ИД - %d", oldUser.Id)
	}
}

func (auth *AuthRepo) GetUser(userName, password string) (internal.User, error) {
	var user internal.User
	hashPass := EncodePassword(password)
	err := auth.db.Select().From(user.TableName()).Where(dbx.HashExp{"username": userName, "password": hashPass}).One(&user)
	if err != nil {
		return user, fmt.Errorf("неверное имя пользователя или пароль.")
	}
	return user, nil
}

func (auth *AuthRepo) GetAllUser() []internal.User {
	var users []internal.User
	q := auth.db.Select("id", "username", "isadmin").From("users")
	q.All(&users)
	return users
}

func (auth *AuthRepo) DeleteUser(userId int) error {
	_, err := auth.db.Delete("users", dbx.HashExp{"id": userId}).Execute()
	if err != nil {
		return fmt.Errorf("ошибка при удалении пользователя с ИД: %d", userId)
	} else {
		return nil
	}
}
