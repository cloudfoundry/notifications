package fakes

import (
	"database/sql"
	"errors"

	"github.com/cloudfoundry-incubator/notifications/models"
)

type DBResult struct{}

func (fake DBResult) LastInsertId() (int64, error) {
	return 0, nil
}

func (fake DBResult) RowsAffected() (int64, error) {
	return 0, nil
}

type Connection struct {
	BeginWasCalled    bool
	CommitWasCalled   bool
	RollbackWasCalled bool
	CommitError       string

	SelectOneCall struct {
		Returns   interface{}
		Errs      []error
		CallCount int
	}

	SelectCall struct {
		Err error
	}

	InsertCall struct {
		Err error
	}

	UpdateCall struct {
		List []interface{}
	}

	GetCall struct {
		Returns interface{}
		Err     error
	}
}

func NewConnection() *Connection {
	return &Connection{}
}

func (conn *Connection) Begin() error {
	conn.BeginWasCalled = true
	return nil
}

func (conn *Connection) Commit() error {
	conn.CommitWasCalled = true
	if conn.CommitError != "" {
		return errors.New(conn.CommitError)
	}
	return nil
}

func (conn *Connection) Rollback() error {
	conn.RollbackWasCalled = true
	return nil
}

func (conn *Connection) Exec(query string, args ...interface{}) (sql.Result, error) {
	return DBResult{}, nil
}

func (conn Connection) Delete(list ...interface{}) (int64, error) {
	return 0, nil
}

func (conn Connection) Insert(list ...interface{}) error {
	return conn.InsertCall.Err
}

func (conn Connection) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return []interface{}{}, conn.SelectCall.Err
}

func (conn *Connection) SelectOne(i interface{}, query string, args ...interface{}) error {
	switch returns := conn.SelectOneCall.Returns.(type) {
	case models.Client:
		*i.(*models.Client) = returns
	case models.Kind:
		*i.(*models.Kind) = returns
	}
	call := conn.SelectOneCall.CallCount
	conn.SelectOneCall.CallCount++
	return conn.SelectOneCall.Errs[call]
}

func (conn *Connection) Update(list ...interface{}) (int64, error) {
	conn.UpdateCall.List = list
	return 0, nil
}

func (conn *Connection) Transaction() models.TransactionInterface {
	return conn
}

func (conn *Connection) Get(i interface{}, keys ...interface{}) (interface{}, error) {
	return i, nil
}
