package fakes

import (
    "database/sql"

    "github.com/cloudfoundry-incubator/notifications/models"
)

type FakeDBResult struct{}

func (fake FakeDBResult) LastInsertId() (int64, error) {
    return 0, nil
}

func (fake FakeDBResult) RowsAffected() (int64, error) {
    return 0, nil
}

type FakeDBConn struct {
    BeginWasCalled    bool
    CommitWasCalled   bool
    RollbackWasCalled bool
}

func (conn *FakeDBConn) Begin() error {
    conn.BeginWasCalled = true
    return nil
}

func (conn *FakeDBConn) Commit() error {
    conn.CommitWasCalled = true
    return nil
}

func (conn *FakeDBConn) Rollback() error {
    conn.RollbackWasCalled = true
    return nil
}

func (conn *FakeDBConn) Exec(query string, args ...interface{}) (sql.Result, error) {
    return FakeDBResult{}, nil
}

func (conn FakeDBConn) Delete(list ...interface{}) (int64, error) {
    return 0, nil
}

func (conn FakeDBConn) Insert(list ...interface{}) error {
    return nil
}

func (conn FakeDBConn) Select(i interface{}, query string, args ...interface{}) ([]interface{}, error) {
    return []interface{}{}, nil
}

func (conn FakeDBConn) SelectOne(i interface{}, query string, args ...interface{}) error {
    return nil
}

func (conn FakeDBConn) Update(list ...interface{}) (int64, error) {
    return 0, nil
}

func (conn FakeDBConn) SelectInt(query string, list ...interface{}) (int64, error) {
    return 0, nil
}

func (conn *FakeDBConn) Transaction() models.TransactionInterface {
    return conn
}
