package room

import (
	"context"
	"fmt"
	"log"
	"os"
	"picker/backend/go/pkg/dynamodbTypes"
	"picker/backend/go/pkg/option"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type CreateRoomRequest struct {
	ID       string   `json:"id" binding:"required,alphanum,min=1,max=100"`
	Options  []string `json:"options" binding:"required,gt=0,lt=200,dive,required,min=1,max=1000"`
	Question string   `json:"question" binding:"required,min=1,max=1500"`
}

type Room struct {
	// DynamoDB
	PK     string `dynamodbav:"PK" json:"-"`
	SK     string `dynamodbav:"SK" json:"-"`
	GSI1PK string `dynamodbav:"GSI1PK" json:"-"`
	GSI1SK string `dynamodbav:"GSI1SK" json:"-"`
	Type   string `dynamodbav:"type" json:"-"`

	// Public
	ID       string           `json:"id"`
	Options  []*option.Option `json:"options"`
	Question string           `json:"question"`

	// Private
	OwnerID string `dynamodbav:"ownerID" json:"-"`
}

type PublicRoom struct {
	// Public
	ID       string           `json:"id"`
	Options  []*option.Option `json:"options"`
	Question string           `json:"question"`
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

func Unmarshal(item map[string]types.AttributeValue) *Room {
	room := &Room{}
	if err := attributevalue.UnmarshalMap(item, room); err != nil {
		panic(err)
	}

	return room
}

func NewRoom(request *CreateRoomRequest, userID string, client *dynamodb.Client) (*Room, error) {
	room := &Room{
		PK:       fmt.Sprintf("ROOM#%s", request.ID),
		SK:       fmt.Sprintf("ROOM#%s", request.ID),
		Type:     dynamodbTypes.Room,
		ID:       request.ID,
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

	var options []*option.Option
	for _, opt := range request.Options {
		options = append(options, option.NewOption(opt, userID, request.ID, client))
	}

	err = option.BatchWriteOptions(options, client)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func GetRoom(id string, client *dynamodb.Client) (*Room, error) {
	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("table")),
		KeyConditionExpression: aws.String("PK = :PK and begins_with(SK, :roomPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
			":roomPrefix": &types.AttributeValueMemberS{Value: "ROOM"},
		},
	})

	var room *Room
	var options []*option.Option = []*option.Option{}

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(context.TODO())

		if err != nil {
			panic(err)
		}

		for _, item := range out.Items {
			itemType := dynamodbTypes.GetType(item)

			switch itemType {
			case dynamodbTypes.Room:
				room = Unmarshal(item)
			case dynamodbTypes.Option:
				options = append(options, option.Unmarshal(item))
			default:
				log.Default().Printf("%s missing", itemType)
			}
		}
	}

	room.Options = options

	return room, nil
}
