package dynamodb

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"

type Client interface {
	Save(table string, item interface{}) error
	GetOne(table string, partitionKey string, bindTo interface{}) error
	GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error
	QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error
	BatchGetWithSort(values map[string]interface{}) error
	QueryExpression(table string, query expression.Expression, pageSize int32, pageNumber int32, bindTo interface{}) error
	QueryGSI(table string, globalIndex string, query expression.Expression, customLimit int32, pageDesired int32, bindTo interface{}) error
}
