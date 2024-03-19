package internal

type User struct {
	Id       int    `json:"id"`          // идентификатор пользователя
	Username string `json:"username"`    // имя пользователя
	Password string `json:"password"`    // пароль пользователя
	Isadmin  bool   `json:"userIsAdmin"` // признак того, что пользователь является администратором
}

func (u *User) TableName() string {
	return "users"
}

type DeleteUserStruct struct {
	Id int `json:"id"`
}
