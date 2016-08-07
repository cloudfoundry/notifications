package db

import (
	"database/sql"

	"gopkg.in/gorp.v1"
)

type TransactionInterface interface {
	ConnectionInterface
	Begin() error
	Commit() error
	Rollback() error
}

type Transaction struct {
	txn  *gorp.Transaction
	conn *Connection
}

func NewTransaction(conn *Connection) TransactionInterface {
	return &Transaction{
		conn: conn,
	}
}

func (transaction *Transaction) GetDbMap() *gorp.DbMap {
	return transaction.conn.GetDbMap()
}

func (transaction *Transaction) Begin() error {
	var err error
	transaction.txn, err = transaction.conn.Begin()
	return err
}

func (transaction *Transaction) Transaction() TransactionInterface {
	return transaction
}

func (transaction *Transaction) Commit() error {
	return transaction.txn.Commit()
}

func (transaction *Transaction) Delete(v ...interface{}) (int64, error) {
	return transaction.txn.Delete(v...)
}

func (transaction *Transaction) Exec(query string, v ...interface{}) (sql.Result, error) {
	return transaction.txn.Exec(query, v...)
}

func (transaction *Transaction) Insert(v ...interface{}) error {
	return transaction.txn.Insert(v...)
}

func (transaction *Transaction) Rollback() error {
	return transaction.txn.Rollback()
}

func (transaction *Transaction) Select(holder interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return transaction.txn.Select(holder, query, args...)
}

func (transaction *Transaction) SelectOne(holder interface{}, query string, args ...interface{}) error {
	return transaction.txn.SelectOne(holder, query, args...)
}

func (transaction *Transaction) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return transaction.txn.Get(i, keys)
}

func (transaction *Transaction) Update(v ...interface{}) (int64, error) {
	return transaction.txn.Update(v...)
}
