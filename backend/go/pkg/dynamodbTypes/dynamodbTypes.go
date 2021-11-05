package dynamodbTypes

import (
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const (
	Room   = "room"
	User   = "user"
	Option = "option"
)

type Simple struct {
	Type string
}

func GetType(item map[string]types.AttributeValue) string {
	base := &Simple{}

	err := attributevalue.UnmarshalMap(item, base)

	if err != nil {
		panic(err)
	}

	return base.Type
}
