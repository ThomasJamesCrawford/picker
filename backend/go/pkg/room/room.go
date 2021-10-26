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
	PK   string `dynamodbav:"PK" json:"-"`
	SK   string `dynamodbav:"SK" json:"-"`
	Type string `dynamodbav:"type" json:"-"`

	ID            string `json:"id"`
	OwnershipHash string `json:"ownershipHash"`
}

type PublicRoom struct {
	ID string `json:"id"`
}

func GetPublicRoom(id string, client *dynamodb.Client) (*PublicRoom, error) {
	room, err := GetRoom(id, client)

	if err != nil {
		return nil, err
	}

	return &PublicRoom{ID: room.ID}, nil
}

func CheckRoomNameAvailable(id string, client *dynamodb.Client) bool {
	_, err := GetRoom(id, client)

	return err != nil
}

func GetRoom(id string, client *dynamodb.Client) (*Room, error) {
	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("table")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
	})

	if err != nil {
		return nil, err
	}

	room := &Room{}

	unmarhsalError := attributevalue.UnmarshalMap(res.Item, &room)

	if unmarhsalError != nil {
		panic(unmarhsalError)
	}

	return room, nil
}
