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
)

type CreateRoomRequest struct {
	ID       string   `json:"id" binding:"required,alphanum,min=1,max=100"`
	Options  []string `json:"options" binding:"required,gt=0,dive,required,min=1,max=1000"`
	Question string   `json:"question" binding:"required,min=1,max=1500"`
}

type Room struct {
	PK     string `dynamodbav:"PK" json:"-"`
	SK     string `dynamodbav:"SK" json:"-"`
	GSI1PK string `dynamodbav:"GSI1PK" json:"-"`
	GSI1SK string `dynamodbav:"GSI1SK" json:"-"`
	Type   string `dynamodbav:"type" json:"-"`

	ID       string   `json:"id"`
	Options  []string `json:"options"`
	Question string   `json:"question"`
	OwnerID  string   `json:"ownerID"`
}

type PublicRoom struct {
	ID       string   `json:"id"`
	Options  []string `json:"options"`
	Question string   `json:"question"`
}

func GetPublicRoom(id string, client *dynamodb.Client) (*PublicRoom, error) {
	room, err := GetRoom(id, client)

	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, nil
	}

	return &PublicRoom{
		ID:       room.ID,
		Options:  room.Options,
		Question: room.Question,
	}, nil
}

func NewRoom(request *CreateRoomRequest, userID string, client *dynamodb.Client) (*Room, error) {
	room := &Room{
		PK:       fmt.Sprintf("ROOM#%s", request.ID),
		SK:       fmt.Sprintf("ROOM#%s", request.ID),
		Type:     dynamodbTypes.Room,
		ID:       request.ID,
		Options:  request.Options,
		Question: request.Question,
		OwnerID:  userID,
		GSI1PK:   fmt.Sprintf("USER#%s", userID),
		GSI1SK:   fmt.Sprintf("ROOM#%s", request.ID),
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
