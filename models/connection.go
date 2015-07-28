package models

import (
	"database/sql"

	"github.com/coopernurse/gorp"
)

type ConnectionInterface interface {
	Delete(...interface{}) (int64, error)
	Insert(...interface{}) error
	Select(interface{}, string, ...interface{}) ([]interface{}, error)
	SelectOne(interface{}, string, ...interface{}) error
	Update(...interface{}) (int64, error)
	Exec(string, ...interface{}) (sql.Result, error)
	Transaction() TransactionInterface
	Get(i interface{}, keys ...interface{}) (interface{}, error)
}

type Connection struct {
	*gorp.DbMap
}

func (conn *Connection) Transaction() TransactionInterface {
	return NewTransaction(conn)
}
