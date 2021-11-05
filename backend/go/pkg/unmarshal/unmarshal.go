package unmarshal

import (
	"log"
	"picker/backend/go/pkg/dynamodbTypes"
	"picker/backend/go/pkg/option"
	"picker/backend/go/pkg/room"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Base struct {
	Type string
}

func Unmarshal(item map[string]types.AttributeValue) interface{} {
	base := &Base{}

	err := attributevalue.UnmarshalMap(item, base)

	if err != nil {
		panic(err)
	}

	var ret interface{}

	switch base.Type {
	case dynamodbTypes.Room:
		room := &room.Room{}
		if err = attributevalue.Unmarshal(item, room); err != nil {
			panic(err)
		}
		ret = room
	case dynamodbTypes.Option:
		option := &option.Option{}
		if err = attributevalue.Unmarshal(item, option); err != nil {
			panic(err)
		}
		ret = option
	default:
		log.Default().Printf("%s doesn't exist in unmarshal.Unmarshal", base.Type)
	}

	return ret
}

func Room(item interface{}) *room.Room {
	return item.(*room.Room)
}

func Option(item interface{}) *option.Option {
	return item.(*option.Option)
}
