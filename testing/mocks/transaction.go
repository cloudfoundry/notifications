package mocks

import "github.com/go-gorp/gorp"

type Transaction struct {
	BeginCall struct {
		WasCalled bool
		Returns   struct {
			Error error
		}
	}

	CommitCall struct {
		WasCalled bool
		Returns   struct {
			Error error
		}
	}

	RollbackCall struct {
		WasCalled bool
		Returns   struct {
			Error error
		}
	}

	GetDbMapCall struct {
		WasCalled bool
		Returns   struct {
			DbMap *gorp.DbMap
		}
	}

	*Connection
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) Begin() error {
	t.BeginCall.WasCalled = true
	return t.BeginCall.Returns.Error
}

func (t *Transaction) Commit() error {
	t.CommitCall.WasCalled = true
	return t.CommitCall.Returns.Error
}

func (t *Transaction) Rollback() error {
	t.RollbackCall.WasCalled = true
	return t.RollbackCall.Returns.Error
}

func (t *Transaction) GetDbMap() *gorp.DbMap {
	t.GetDbMapCall.WasCalled = true
	return t.GetDbMapCall.Returns.DbMap
}
