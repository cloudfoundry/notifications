package fakes

import (
    "database/sql"

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
    return nil
}

func (conn DBConn) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
    return []interface{}{}, nil
}

func (conn DBConn) SelectOne(i interface{}, query string, args ...interface{}) error {
    return nil
}

func (conn DBConn) Update(list ...interface{}) (int64, error) {
    return 0, nil
}

func (conn *DBConn) Transaction() models.TransactionInterface {
    return conn
}
