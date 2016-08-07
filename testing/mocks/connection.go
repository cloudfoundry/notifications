package mocks

import (
	"database/sql"

	"github.com/cloudfoundry-incubator/notifications/db"
	"gopkg.in/gorp.v1"
)

type Connection struct {
	DeleteCall struct {
		Receives struct {
			List []interface{}
		}
		Returns struct {
			Count int64
			Error error
		}
	}

	ExecCall struct {
		Receives struct {
			Query  string
			Params []interface{}
		}
		Returns struct {
			Result sql.Result
			Error  error
		}
	}

	GetCall struct {
		Receives struct {
			Object interface{}
			Keys   []interface{}
		}
		Returns struct {
			Object interface{}
			Error  error
		}
	}

	InsertCall struct {
		Receives struct {
			List []interface{}
		}
		Returns struct {
			Error error
		}
	}

	SelectCall struct {
		Receives struct {
			ResultType interface{}
			Query      string
			Params     []interface{}
		}
		Returns struct {
			Results []interface{}
			Error   error
		}
	}

	SelectOneCall struct {
		Receives struct {
			Result interface{}
			Query  string
			Params []interface{}
		}
		Returns struct {
			Error error
		}
	}

	TransactionCall struct {
		Returns struct {
			Transaction db.TransactionInterface
		}
	}

	UpdateCall struct {
		Receives struct {
			List []interface{}
		}
		Returns struct {
			Count int64
			Error error
		}
	}

	GetDbMapCall struct {
		WasCalled bool
		Returns   struct {
			DbMap *gorp.DbMap
		}
	}
}

func NewConnection() *Connection {
	return &Connection{}
}

func (c *Connection) Delete(list ...interface{}) (int64, error) {
	c.DeleteCall.Receives.List = list

	return c.DeleteCall.Returns.Count, c.DeleteCall.Returns.Error
}

func (c *Connection) Exec(query string, params ...interface{}) (sql.Result, error) {
	c.ExecCall.Receives.Query = query
	c.ExecCall.Receives.Params = params

	return c.ExecCall.Returns.Result, c.ExecCall.Returns.Error
}

func (c *Connection) Get(object interface{}, keys ...interface{}) (interface{}, error) {
	c.GetCall.Receives.Object = object
	c.GetCall.Receives.Keys = keys

	return c.GetCall.Returns.Object, c.GetCall.Returns.Error
}

func (c *Connection) Insert(list ...interface{}) error {
	c.InsertCall.Receives.List = list

	return c.InsertCall.Returns.Error
}

func (c *Connection) Select(resultType interface{}, query string, params ...interface{}) ([]interface{}, error) {
	c.SelectCall.Receives.ResultType = resultType
	c.SelectCall.Receives.Query = query
	c.SelectCall.Receives.Params = params

	return c.SelectCall.Returns.Results, c.SelectCall.Returns.Error
}

func (c *Connection) SelectOne(result interface{}, query string, params ...interface{}) error {
	c.SelectOneCall.Receives.Result = result
	c.SelectOneCall.Receives.Query = query
	c.SelectOneCall.Receives.Params = params

	return c.SelectOneCall.Returns.Error
}

func (c *Connection) Transaction() db.TransactionInterface {
	return c.TransactionCall.Returns.Transaction
}

func (c *Connection) Update(list ...interface{}) (int64, error) {
	c.UpdateCall.Receives.List = list

	return c.UpdateCall.Returns.Count, c.UpdateCall.Returns.Error
}

func (c *Connection) GetDbMap() *gorp.DbMap {
	c.GetDbMapCall.WasCalled = true
	return c.GetDbMapCall.Returns.DbMap
}
