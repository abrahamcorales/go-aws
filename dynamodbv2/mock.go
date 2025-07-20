package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/stretchr/testify/mock"
)

type DynamoMock struct {
	mock.Mock
}

func (mock *DynamoMock) Save(table string, item interface{}) error {
	args := mock.Called(item)
	return args.Error(0)
}

func (mock *DynamoMock) GetOne(table string, partitionKey string, bindTo interface{}) error {
	args := mock.Called(partitionKey, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error {
	args := mock.Called(partitionKey, sortKey, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error {
	args := mock.Called(partitionKey, limit, bindTo)
	return args.Error(0)
}

func (mock *DynamoMock) BatchGetWithSort(values map[string]interface{}) error {
	ret := mock.Called(values)

	var r0 error
	if rf, ok := ret.Get(0).(func(map[string]interface{}) error); ok {
		r0 = rf(values)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryMultiple provides a mock function with given fields: table, partitionKey, limit, bindTo
func (_m *DynamoMock) QueryMultiple(table string, partitionKey string, limit int32, bindTo interface{}) error {
	ret := _m.Called(table, partitionKey, limit, bindTo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, int32, interface{}) error); ok {
		r0 = rf(table, partitionKey, limit, bindTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryExpression provides a mock function with given fields: table, query, pageSize, pageNumber, bindTo
func (_m *DynamoMock) QueryExpression(table string, query expression.Expression, pageSize int32, pageNumber int32, bindTo interface{}) error {
	ret := _m.Called(table, query, pageSize, pageNumber, bindTo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, expression.Expression, int32, int32, interface{}) error); ok {
		r0 = rf(table, query, pageSize, pageNumber, bindTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryGSI provides a mock function with given fields: table, globalIndex, query, customLimit, pageDesired, bindTo
func (_m *DynamoMock) QueryGSI(table string, globalIndex string, query expression.Expression, customLimit int32, pageDesired int32, bindTo interface{}) error {
	ret := _m.Called(table, globalIndex, query, customLimit, pageDesired, bindTo)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, expression.Expression, int32, int32, interface{}) error); ok {
		r0 = rf(table, globalIndex, query, customLimit, pageDesired, bindTo)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
