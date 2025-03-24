// Code generated by mockery v2.53.3. DO NOT EDIT.

package pgrx

import (
	context "context"

	pgconn "github.com/jackc/pgx/v5/pgconn"
	mock "github.com/stretchr/testify/mock"
)

// MockDB is an autogenerated mock type for the DB type
type MockDB struct {
	mock.Mock
}

type MockDB_Expecter struct {
	mock *mock.Mock
}

func (_m *MockDB) EXPECT() *MockDB_Expecter {
	return &MockDB_Expecter{mock: &_m.Mock}
}

// Exec provides a mock function with given fields: ctx, sql, arguments
func (_m *MockDB) Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error) {
	var _ca []interface{}
	_ca = append(_ca, ctx, sql)
	_ca = append(_ca, arguments...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Exec")
	}

	var r0 pgconn.CommandTag
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)); ok {
		return rf(ctx, sql, arguments...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, ...interface{}) pgconn.CommandTag); ok {
		r0 = rf(ctx, sql, arguments...)
	} else {
		r0 = ret.Get(0).(pgconn.CommandTag)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, ...interface{}) error); ok {
		r1 = rf(ctx, sql, arguments...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockDB_Exec_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Exec'
type MockDB_Exec_Call struct {
	*mock.Call
}

// Exec is a helper method to define mock.On call
//   - ctx context.Context
//   - sql string
//   - arguments ...interface{}
func (_e *MockDB_Expecter) Exec(ctx interface{}, sql interface{}, arguments ...interface{}) *MockDB_Exec_Call {
	return &MockDB_Exec_Call{Call: _e.mock.On("Exec",
		append([]interface{}{ctx, sql}, arguments...)...)}
}

func (_c *MockDB_Exec_Call) Run(run func(ctx context.Context, sql string, arguments ...interface{})) *MockDB_Exec_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]interface{}, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(interface{})
			}
		}
		run(args[0].(context.Context), args[1].(string), variadicArgs...)
	})
	return _c
}

func (_c *MockDB_Exec_Call) Return(_a0 pgconn.CommandTag, _a1 error) *MockDB_Exec_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockDB_Exec_Call) RunAndReturn(run func(context.Context, string, ...interface{}) (pgconn.CommandTag, error)) *MockDB_Exec_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockDB creates a new instance of MockDB. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockDB(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockDB {
	mock := &MockDB{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
