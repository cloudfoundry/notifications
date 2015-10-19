package mocks

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
