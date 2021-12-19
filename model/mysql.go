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
func (m *MySQLHandler) ChangeUserAuth(userAuth, userId int) (*User, error) {
	return nil, nil
}
func (m *MySQLHandler) DeleteUser(userAuth, userId int) (bool, error) {
	return false, nil
}

func newMySQLHandler() DBHandler {
	database, err := sql.Open("mysql", "root:root@tcp(127.0.0.1:3306)/sgs")
	if err != nil {
		panic(err)
	}
	fmt.Print(database)
	// statement, _ := database.Prepare(
	// 	`CREATE TABLE IF NOT EXISTS users (
	// 		id int primary key auto_increment,
	// 		name varchar(32) not null,
	// 		phone varchar(12) not null,
	// 		email varchar(64) not null,
	// 		authority tinyint(1) not null,
	// 		password varchar(128) not null,
	// 		created_at datetime DEFAULT CURRENT_TIMESTAMP not null
	// 	)ENGINE=INNODB;
	// 	CREATE INDEX IF NOT EXISTS idIndexOnUsers ON todos (
	// 		id ASC
	// 	);`)
	// statement.Exec()
	return &MySQLHandler{db: database}
}
