package room

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Room struct {
	PK   string `dynamodbav:"PK"`
	SK   string `dynamodbav:"SK"`
	Type string `dynamodbav:"type"`

	Name       string `json:"name"`
	SecretName string `json:"secretName"`
}

func GetRoom(id string, client *dynamodb.Client) *Room {
	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("table")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
	})

	if err != nil {
		panic(err)
	}

	room := &Room{}

	unmarhsalError := attributevalue.UnmarshalMap(res.Item, &room)

	if unmarhsalError != nil {
		panic(unmarhsalError)
	}

	return room
}
