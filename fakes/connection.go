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

type DBConn struct {
	BeginWasCalled    bool
	CommitWasCalled   bool
	RollbackWasCalled bool
	CommitError       string

	SelectOneCall struct {
		Returns   interface{}
		Errs      []error
		CallCount int
	}

	InsertCall struct {
		Err error
	}

	UpdateCall struct {
		List []interface{}
	}
}

func NewDBConn() *DBConn {
	return &DBConn{}
}

func (conn *DBConn) Begin() error {
	conn.BeginWasCalled = true
	return nil
}

func (conn *DBConn) Commit() error {
	conn.CommitWasCalled = true
	if conn.CommitError != "" {
		return errors.New(conn.CommitError)
	}
	return nil
}

func (conn *DBConn) Rollback() error {
	conn.RollbackWasCalled = true
	return nil
}

func (conn *DBConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return DBResult{}, nil
}

func (conn DBConn) Delete(list ...interface{}) (int64, error) {
	return 0, nil
}

func (conn DBConn) Insert(list ...interface{}) error {
	return conn.InsertCall.Err
}

func (conn DBConn) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
	return []interface{}{}, nil
}

func (conn *DBConn) SelectOne(i interface{}, query string, args ...interface{}) error {
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

func (conn *DBConn) Update(list ...interface{}) (int64, error) {
	conn.UpdateCall.List = list
	return 0, nil
}

func (conn *DBConn) Transaction() models.TransactionInterface {
	return conn
}
