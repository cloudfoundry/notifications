package db

import (
	"database/sql"

	"gopkg.in/gorp.v1"
)

type ConnectionInterface interface {
	Transaction() TransactionInterface
	GetDbMap() *gorp.DbMap
	Delete(...interface{}) (int64, error)
	Insert(...interface{}) error
	Select(interface{}, string, ...interface{}) ([]interface{}, error)
	SelectOne(interface{}, string, ...interface{}) error
	Update(...interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Get(i interface{}, keys ...interface{}) (interface{}, error)
}

type Connection struct {
	*gorp.DbMap
}

func (conn *Connection) Transaction() TransactionInterface {
	return NewTransaction(conn)
}

func (conn *Connection) GetDbMap() *gorp.DbMap {
	return conn.DbMap
}
