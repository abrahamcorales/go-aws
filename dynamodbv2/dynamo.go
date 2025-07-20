package dynamodb

import (
	"context"
	"errors"
	"fmt"

	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var (
	// ErrNotFound is returned when no items could be found in Get or OldValue and similar operations.
	ErrNotFound = errors.New("dynamo: no item found")
	Tagkey      = "dynamo"
	tables      = make(map[string]DynamoTable)
)

type Implementation struct {
	client            *dynamodb.Client
	partitionKeyField string
	sortKeyField      *string
	DynamoTables      map[string]DynamoTable
}

type DynamoTable struct {
	TableName         string `json:"table_name"`
	PartitionKeyField string `json:"primary_key_field"`
	SortKeyField      string `json:"sort_key_field"`
	MaxPageSize       int32  `json:"max_page_size"`
	GlobalIndex       string `json:"global_index"`
}

type funcTable func(i *Implementation)

func WithTable(arg DynamoTable) funcTable {
	return func(i *Implementation) {
		tables[arg.TableName] = DynamoTable{
			TableName:         arg.TableName,
			PartitionKeyField: arg.PartitionKeyField,
			SortKeyField:      arg.SortKeyField,
			MaxPageSize:       arg.MaxPageSize,
		}
		i.DynamoTables = tables
	}
}

func NewDynamoClientv2(awsConfig aws.Config, funcTableArray ...funcTable) Client {
	var i Implementation
	db := dynamodb.NewFromConfig(awsConfig)
	for _, ft := range funcTableArray {
		ft(&i)
	}
	i.client = db
	return &i
}

func (i *Implementation) Save(table string, values interface{}) error {
	log.Printf("[DynamoDB] executing put query")

	item, err := attributevalue.MarshalMapWithOptions(values, func(h *attributevalue.EncoderOptions) {
		h.TagKey = Tagkey
	})

	log.Printf("[DynamoDB] item to save: %s", item)
	if err != nil {
		panic(fmt.Sprintf("failed to DynamoDB marshal Record, %v", err))
	}

	_, err = i.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(i.DynamoTables[table].TableName),
		Item:      item,
	})
	return err
}

func (i *Implementation) getItem(table string, key map[string]types.AttributeValue, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query")

	out, err := i.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(i.DynamoTables[table].TableName),
		Key:       key,
	})
	if err != nil {
		return err
	}
	if out.Item == nil {
		return ErrNotFound
	}
	err = attributevalue.UnmarshalMapWithOptions(out.Item, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})

	return err
}

func (i *Implementation) getItemQuery(table string, key string, limit int32, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query")
	keyEx := expression.Key(i.DynamoTables[table].PartitionKeyField).Equal(expression.Value(key))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return err
	}
	out, err := i.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:                 aws.String(i.DynamoTables[table].TableName),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
		Limit:                     aws.Int32(limit),
	})
	if out != nil {
		if len(out.Items) == 0 {
			return ErrNotFound
		}
	}
	if err != nil {
		return err
	}

	err = attributevalue.UnmarshalListOfMapsWithOptions(out.Items, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})
	return err
}

func (i *Implementation) ItemQueryExpression(table string, globalIndex string, query expression.Expression, pageSize int32, pageNumber int32, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query") // Log indicating that a DynamoDB query is being executed

	// Get MaxPageSize from table configuration or custom max page size
	maxPageSize := getLimitPageSize(i.DynamoTables[table].MaxPageSize, pageSize)

	var (
		lastEvaluatedKey     map[string]types.AttributeValue   // Last evaluated key for pagination
		itemsTotal           []map[string]types.AttributeValue // Total items retrieved
		totalConsumeCapacity float64                           // Total consumed capacity
		page                 int                               // Current page number
		count                int32                             // Count of items retrieved
		err                  error                             // Error
	)

	// Build the query input
	queryInput := buildQueryInput(i.DynamoTables[table].TableName, globalIndex, query, lastEvaluatedKey)

	// Apply max limit item if pageSize is set
	applyLimits(&queryInput.Limit, maxPageSize)

	for {
		// Update ExclusiveStartKey with the next items to be fetched
		queryInput.ExclusiveStartKey = lastEvaluatedKey
		output, err := i.client.Query(context.TODO(), queryInput)
		if err != nil {
			return fmt.Errorf("failed to query items: %w", err) // Return error if the query fails
		}

		totalConsumeCapacity += *output.ConsumedCapacity.CapacityUnits // Add consumed capacity
		page++                                                         // Increment page number
		// Append items to total items
		itemsTotal = append(itemsTotal, output.Items...)
		count += output.Count // Increment the count of items retrieved

		// Check if the desired page has been setup and reached or there are no more items
		if pageNumber > 0 && (hasReachDesiredPage(page, pageNumber) || output.LastEvaluatedKey == nil) {
			itemsTotal = output.Items
			logQueryStatus(query.Values(), totalConsumeCapacity, page, count) // Log the query status
			break
		}
		// Check if the limit of items has been reached or there are no more items
		if hasReachedLimit(count, maxPageSize, pageNumber) || output.LastEvaluatedKey == nil {
			logQueryStatus(query.Values(), totalConsumeCapacity, page, count) // Log the query status
			break
		}
		// Update lastEvaluatedKey
		lastEvaluatedKey = output.LastEvaluatedKey

		// Update the query limit
		updateQueryLimit(&queryInput.Limit, &maxPageSize, count)
	}

	// Deserialize the list of attribute maps into bindTo
	err = attributevalue.UnmarshalListOfMapsWithOptions(itemsTotal, &bindTo, func(options *attributevalue.DecoderOptions) {
		options.TagKey = Tagkey
	})
	return err // Return error if any
}

func (i *Implementation) BatchGetItem(key map[string]types.KeysAndAttributes, values map[string]interface{}) error {
	log.Printf("[DynamoDB] batch get query")
	out, err := i.client.BatchGetItem(context.TODO(), &dynamodb.BatchGetItemInput{
		RequestItems: key,
	})
	if err != nil {
		return err
	}
	var bindList []interface{}
	for i, o := range out.Responses {
		for t, v := range values {
			if i == t {
				v, _ := v.([]interface{})
				bindList := append(bindList, v[2])
				err = attributevalue.UnmarshalListOfMapsWithOptions(o, &bindList, func(options *attributevalue.DecoderOptions) {
					options.TagKey = Tagkey
				})
			}
		}
	}
	return err
}

func (i *Implementation) GetOne(table string, partitionKey string, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query")
	return i.getItem(table,
		map[string]types.AttributeValue{
			i.DynamoTables[table].PartitionKeyField: &types.AttributeValueMemberS{Value: partitionKey},
		},
		bindTo)
}

func (i *Implementation) GetOneWithSort(table string, partitionKey string, sortKey string, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query with sortkey [pk:%s][sk:%s]", partitionKey, sortKey)

	return i.getItem(table,
		map[string]types.AttributeValue{
			i.partitionKeyField: &types.AttributeValueMemberS{Value: partitionKey},
			*i.sortKeyField:     &types.AttributeValueMemberS{Value: sortKey},
		},
		bindTo)
}

func (i *Implementation) QueryOne(table string, partitionKey string, limit int32, bindTo interface{}) error {
	log.Printf("[DynamoDB] executing get query with [pk:%s]", partitionKey)

	return i.getItemQuery(table, partitionKey, limit, bindTo)

}

func (i *Implementation) BatchGetWithSort(values map[string]interface{}) error {
	log.Printf("[DynamoDB] executing BatchGet query")
	batchkeys := make(map[string]types.KeysAndAttributes)
	for t, v := range values {
		v, _ := v.([]interface{})
		batchkeys[t] = types.KeysAndAttributes{
			Keys: []map[string]types.AttributeValue{
				{
					i.DynamoTables[t].PartitionKeyField: &types.AttributeValueMemberS{Value: v[0].(string)},
					i.DynamoTables[t].SortKeyField:      &types.AttributeValueMemberS{Value: v[1].(string)},
				},
			},
		}

	}

	return i.BatchGetItem(batchkeys, values)
}

// QueryExpression  returns multiple items by using a query expression
func (i *Implementation) QueryExpression(table string, query expression.Expression, pageSize int32, pageNumber int32, bindTo interface{}) error {

	return i.ItemQueryExpression(table, "", query, pageSize, pageNumber, bindTo)
}

func applyLimits(limit **int32, limitMaxItems int32) {
	if limitMaxItems > 0 {
		*limit = aws.Int32(limitMaxItems)
	}
}

func hasReachedLimit(count, limitMaxItems int32, pageNumber int32) bool {
	return (limitMaxItems > 0 && count >= limitMaxItems && pageNumber <= 0)
}

func hasReachDesiredPage(pageNumber int, pageDesired int32) bool {
	return (int32(pageNumber) >= pageDesired)
}

func updateQueryLimit(limit **int32, limitMaxItems *int32, count int32) {
	if *limitMaxItems > count {
		*limitMaxItems -= count
		*limit = aws.Int32(*limitMaxItems)
	}
}

func (i *Implementation) QueryGSI(table string, globalIndex string, query expression.Expression, pageSize int32, pageNumber int32, bindTo interface{}) error {
	return i.ItemQueryExpression(table, globalIndex, query, pageSize, pageNumber, bindTo)
}

func buildQueryInput(tableName, globalIndex string, query expression.Expression, startKey map[string]types.AttributeValue) *dynamodb.QueryInput {
	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(tableName),
		ExpressionAttributeNames:  query.Names(),
		ExpressionAttributeValues: query.Values(),
		KeyConditionExpression:    query.KeyCondition(),
		ExclusiveStartKey:         startKey,
		ReturnConsumedCapacity:    types.ReturnConsumedCapacityTotal,
	}

	if globalIndex != "" {
		queryInput.IndexName = aws.String(globalIndex)
	}
	if query.Projection() != nil {
		queryInput.ProjectionExpression = query.Projection()
	}
	if query.Filter() != nil {
		queryInput.FilterExpression = query.Filter()
	}
	return queryInput
}

func logQueryStatus(pk map[string]types.AttributeValue, totalConsumeCapacity float64, page int, count int32) {
	log.Printf("[DynamoDB] pk: %s | consume Capacity: %v | total Page: %d | count: %d", pk, totalConsumeCapacity, page, count)
}

func getLimitPageSize(defaultLimit, pageSize int32) int32 {
	if pageSize > 0 {
		return pageSize
	}
	return defaultLimit
}
