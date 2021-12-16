package model

import "database/sql"

type MysqlHandler struct {
	db *sql.DB
}

func (m *MysqlHandler) Close() {
	m.db.Close()
}
