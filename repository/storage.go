package repository

import (
	"fmt"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"log/slog"
	"net/http"
	"time"
)

type Storage struct {
	DB *dbx.DB
}

// New возвращает соединение с БД
func New(storagePath string) (*Storage, error) {
	db, err := dbx.MustOpen("postgres", storagePath)
	if err != nil {
		return nil, err
	}
	return &Storage{DB: db}, nil
}

// Init инициализирует БД
func (storage *Storage) Init() error {
	//region Создаем таблицу Пользователей, индекс и пользователя администратора
	queryText := `CREATE TABLE IF NOT EXISTS users (
								id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
								userName varchar(20) NOT NULL,
								password varchar(20) NOT NULL,
								isAdmin boolean	 NOT NULL
								)`
	q := storage.DB.NewQuery(queryText)
	_, err := q.Execute()
	if err != nil {
		return fmt.Errorf("create table 'Users' complete with error: %s", err.Error())
	}

	//проверяем существует ли администратор, если нет - создаем.
	queryText = `Select count(*) as count from USERS where username ='admin'`
	q = storage.DB.NewQuery(queryText)
	resultRows, err := q.Rows()
	if err != nil {
		return fmt.Errorf("check Admin user complete with error: %s", err)
	}
	count := 0
	for resultRows.Next() {
		resultRows.Scan(&count)
	}

	if count == 0 {
		queryText = "INSERT INTO USERS(username, password, isAdmin) VALUES ( 'admin', " + "'" + EncodePassword("123") + "', true )"
		q = storage.DB.NewQuery(queryText)
		_, err = q.Execute()
		if err != nil {
			return fmt.Errorf("create Admin user complete with error: %s", err)
		}
	}

	//endregion

	//region Создаем таблицу Заданий
	queryText = `CREATE TABLE IF NOT EXISTS quests (
								id integer GENERATED BY DEFAULT AS IDENTITY,
								questName varchar(200) NOT NULL,
    							PRIMARY KEY (questName)
								)`
	q = storage.DB.NewQuery(queryText)
	_, err = q.Execute()
	if err != nil {
		return fmt.Errorf("create table 'quests' complete with error: %s", err.Error())
	}
	//endregion

	//region Создаем таблицу questSteps
	queryText = `CREATE TABLE IF NOT EXISTS questSteps (
								id integer PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
								questID integer NOT NULL,
								stepName varchar(200) NOT NULL,
								bonus integer,
    							isMulti bool NOT NULL
								)`
	q = storage.DB.NewQuery(queryText)
	_, err = q.Execute()
	if err != nil {
		return fmt.Errorf("create questSteps 'films' complete with error: %s", err.Error())
	}

	//region Создаем таблицу history
	queryText = `CREATE TABLE IF NOT EXISTS history (
								stepId integer NOT NULL,
								userId integer NOT NULL
								)`
	q = storage.DB.NewQuery(queryText)
	_, err = q.Execute()
	if err != nil {
		return fmt.Errorf("create table 'history' complete with error: %s", err.Error())
	}
	//endregion

	return nil
}

//region Системные методы

func RequestTolog(r *http.Request, logger *slog.Logger) {
	username, _, ok := r.BasicAuth()
	if !ok {
		username = ""
	}
	logger.Debug(
		"incoming request",
		"url", r.URL.Path,
		"method", r.Method,
		"user", username,
		"time", time.Now().Format("02-01-2006 15:04:05"),
	)
}

//endregion Системные методы
