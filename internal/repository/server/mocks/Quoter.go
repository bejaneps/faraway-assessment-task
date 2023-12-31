// Code generated by mockery v2.23.4. DO NOT EDIT.

package server

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Quoter is an autogenerated mock type for the Quoter type
type Quoter struct {
	mock.Mock
}

type Quoter_Expecter struct {
	mock *mock.Mock
}

func (_m *Quoter) EXPECT() *Quoter_Expecter {
	return &Quoter_Expecter{mock: &_m.Mock}
}

// Quote provides a mock function with given fields: ctx
func (_m *Quoter) Quote(ctx context.Context) (string, error) {
	ret := _m.Called(ctx)

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (string, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) string); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Quoter_Quote_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Quote'
type Quoter_Quote_Call struct {
	*mock.Call
}

// Quote is a helper method to define mock.On call
//   - ctx context.Context
func (_e *Quoter_Expecter) Quote(ctx interface{}) *Quoter_Quote_Call {
	return &Quoter_Quote_Call{Call: _e.mock.On("Quote", ctx)}
}

func (_c *Quoter_Quote_Call) Run(run func(ctx context.Context)) *Quoter_Quote_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Quoter_Quote_Call) Return(quote string, err error) *Quoter_Quote_Call {
	_c.Call.Return(quote, err)
	return _c
}

func (_c *Quoter_Quote_Call) RunAndReturn(run func(context.Context) (string, error)) *Quoter_Quote_Call {
	_c.Call.Return(run)
	return _c
}

// NewQuoter creates a new instance of Quoter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewQuoter(t interface {
	mock.TestingT
	Cleanup(func())
}) *Quoter {
	mock := &Quoter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
