package model

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type MySQLHandler struct {
	db *sql.DB
}

func (m *MySQLHandler) Close() {
	m.db.Close()
}

func (m *MySQLHandler) AddUser(name, email, password, phone string) error {
	stmt, err := m.db.Prepare("INSERT INTO users (name, password, email, phone) VALUES(?, ?, ?, ?)")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(name, password, email, phone)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQLHandler) GetUsers() ([]*User, error) {
	users := []*User{}
	rows, err := m.db.Query("SELECT id, name, authority, password, phone, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		rows.Scan(&user.ID, &user.Name, &user.Authority, &user.Password, &user.Phone, &user.Email)
		users = append(users, &user)
	}
	return users, nil
}
func (m *MySQLHandler) GetUser(email string) (*User, error) {
	user := &User{}
	err := m.db.QueryRow("SELECT id, name, authority, password, phone, email FROM users WHERE email=?", email).Scan(&user.ID, &user.Name, &user.Authority, &user.Password, &user.Phone, &user.Email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (m *MySQLHandler) ChangeUserAuth(userAuth string, userId int) error {
	stmt, err := m.db.Prepare("UPDATE users SET authority=? WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userAuth, userId)
	if err != nil {
		return err
	}
	return nil
}
func (m *MySQLHandler) DeleteUser(userId int) error {
	stmt, err := m.db.Prepare("DELETE FROM users WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(userId)
	if err != nil {
		return err
	}
	return nil
}

func newMySQLHandler() DBHandler {
	database, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/sgs")
	if err != nil {
		panic(err)
	}
	fmt.Print(database)
	return &MySQLHandler{db: database}
}
