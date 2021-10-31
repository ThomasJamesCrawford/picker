package room

import (
	"context"
	"fmt"
	"os"
	"picker/backend/go/pkg/dynamodbTypes"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/twinj/uuid"
)

type CreateRoomRequest struct {
	ID      string   `json:"id,omitempty"`
	Options []string `json:"options"`
}

type Room struct {
	PK     string `dynamodbav:"PK" json:"-"`
	SK     string `dynamodbav:"SK" json:"-"`
	GSI1PK string `dynamodbav:"GSI1PK" json:"-"`
	GSI1SK string `dynamodbav:"GSI1SK" json:"-"`
	Type   string `dynamodbav:"type" json:"-"`

	ID            string   `json:"id"`
	OwnershipHash string   `json:"ownershipHash"`
	Options       []string `json:"options"`
}

type PublicRoom struct {
	ID string `json:"id"`
}

func GetPublicRoom(id string, client *dynamodb.Client) (*PublicRoom, error) {
	room, err := GetRoom(id, client)

	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, nil
	}

	return &PublicRoom{ID: room.ID}, nil
}

func NewRoom(request *CreateRoomRequest, client *dynamodb.Client) (*Room, error) {
	ownershipHash := uuid.NewV4().String()

	room := &Room{
		PK:            fmt.Sprintf("ROOM#%s", request.ID),
		SK:            fmt.Sprintf("ROOM#%s", request.ID),
		Type:          dynamodbTypes.Room,
		ID:            request.ID,
		OwnershipHash: ownershipHash,
		Options:       request.Options,
		GSI1PK:        fmt.Sprintf("ROOM#%s", ownershipHash),
		GSI1SK:        fmt.Sprintf("ROOM#%s", ownershipHash),
	}

	marshalledRoom, marshallErr := attributevalue.MarshalMap(room)

	if marshallErr != nil {
		panic(marshallErr)
	}

	_, err := client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("table")),
		Item:                marshalledRoom,
		ConditionExpression: aws.String("attribute_not_exists(PK) and attribute_not_exists(SK)"),
	})

	if err != nil {
		return nil, err
	}

	return room, nil
}

func GetRoom(id string, client *dynamodb.Client) (*Room, error) {
	res, err := client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv("table")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
		},
	})

	if res.Item == nil {
		return nil, nil
	}

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
