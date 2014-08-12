package models

import "github.com/coopernurse/gorp"

type TransactionInterface interface {
    Begin() error
    Commit() error
    Rollback() error
    ConnectionInterface
}

type Transaction struct {
    *gorp.Transaction
}

func NewTransaction() *Transaction {
    return &Transaction{}
}

func (transaction *Transaction) Begin() error {
    var err error
    transaction.Transaction, err = Database().Connection.Begin()
    return err
}
