package room

import (
	"context"
	"fmt"
	"log"
	"os"
	"picker/backend/go/pkg/dynamodbTypes"
	"picker/backend/go/pkg/option"
	"time"

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

type UpdateRoomRequest struct {
	Question string `json:"question" binding:"required,min=1,max=1500"`
}

type Room struct {
	// DynamoDB
	PK     string `dynamodbav:"PK" json:"-"`
	SK     string `dynamodbav:"SK" json:"-"`
	GSI1PK string `dynamodbav:"GSI1PK" json:"-"`
	GSI1SK string `dynamodbav:"GSI1SK" json:"-"`
	Type   string `dynamodbav:"type" json:"-"`

	// Public
	ID       string          `json:"id"`
	Options  []option.Option `json:"options"`
	Question string          `json:"question"`

	// Private
	OwnerID   string    `dynamodbav:"ownerID" json:"-"`
	CreatedAt time.Time `dynamodbav:"createdAt" json:"-"`
}

type PublicRoom struct {
	// Public
	ID        string                `json:"id"`
	Options   []option.PublicOption `json:"options"`
	Question  string                `json:"question"`
	OwnedByMe bool                  `json:"ownedByMe"`
}

func (room Room) getPublic(userID string) PublicRoom {
	publicOptions := option.MapToPublic(room.Options, userID)

	return PublicRoom{
		ID:        room.ID,
		Options:   publicOptions,
		Question:  room.Question,
		OwnedByMe: room.OwnerID == userID,
	}
}

func GetPublicRoom(id string, client *dynamodb.Client, userID string) (*PublicRoom, error) {
	room, err := GetRoom(id, client, userID)

	if err != nil {
		return nil, err
	}

	if room == nil {
		return nil, nil
	}
	publicRoom := room.getPublic(userID)

	return &publicRoom, nil

}

func Unmarshal(item map[string]types.AttributeValue) Room {
	room := &Room{}
	if err := attributevalue.UnmarshalMap(item, room); err != nil {
		panic(err)
	}

	return *room
}

func NewRoom(request *CreateRoomRequest, userID string, client *dynamodb.Client) (*Room, error) {
	createdAt := time.Now().UTC()

	room := &Room{
		PK:        fmt.Sprintf("ROOM#%s", request.ID),
		SK:        fmt.Sprintf("ROOM#%s", request.ID),
		Type:      dynamodbTypes.Room,
		ID:        request.ID,
		Question:  request.Question,
		OwnerID:   userID,
		CreatedAt: createdAt,
		GSI1PK:    fmt.Sprintf("USER#%s", userID),
		GSI1SK:    fmt.Sprintf("ROOM#%s#%s", createdAt.Format(time.RFC3339), request.ID),
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
		newOpt := option.NewOption(opt, userID, request.ID)
		options = append(options, &newOpt)
	}

	err = option.BatchWriteOptions(options, client)

	if err != nil {
		return nil, err
	}

	return room, nil
}

func GetRoom(id string, client *dynamodb.Client, userID string) (*Room, error) {
	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName:              aws.String(os.Getenv("table")),
		KeyConditionExpression: aws.String("PK = :PK and begins_with(SK, :roomPrefix)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":PK":         &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", id)},
			":roomPrefix": &types.AttributeValueMemberS{Value: "ROOM"},
		},
	})

	var room *Room
	var options []option.Option = []option.Option{}

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(context.TODO())

		if err != nil {
			panic(err)
		}

		for _, item := range out.Items {
			itemType := dynamodbTypes.GetType(item)

			switch itemType {
			case dynamodbTypes.Room:
				res := Unmarshal(item)
				room = &res
			case dynamodbTypes.Option:
				options = append(options, option.Unmarshal(item))
			default:
				log.Default().Printf("%s missing", itemType)
			}
		}
	}

	if room == nil {
		return nil, nil
	}

	room.Options = options

	return room, nil
}

func RoomsForUser(userID string, client *dynamodb.Client) (*[]Room, error) {
	paginator := dynamodb.NewQueryPaginator(client, &dynamodb.QueryInput{
		TableName: aws.String(os.Getenv("table")),
		IndexName: aws.String("GSI1"),
		// Get the most recent first
		ScanIndexForward:       aws.Bool(false),
		KeyConditionExpression: aws.String("GSI1PK = :GSI1PK and begins_with(GSI1SK, :room)"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":GSI1PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("USER#%s", userID)},
			":room":   &types.AttributeValueMemberS{Value: "ROOM#"},
		},
	})

	var rooms []Room = []Room{}

	for paginator.HasMorePages() {
		out, err := paginator.NextPage(context.TODO())

		if err != nil {
			panic(err)
		}

		for _, item := range out.Items {
			rooms = append(rooms, Unmarshal(item))
		}
	}

	return &rooms, nil
}

func Update(userID string, roomID string, request UpdateRoomRequest, client *dynamodb.Client) (*Room, error) {
	res, err := client.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
		TableName: aws.String(os.Getenv("table")),
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", roomID)},
			"SK": &types.AttributeValueMemberS{Value: fmt.Sprintf("ROOM#%s", roomID)},
		},
		UpdateExpression: aws.String("set question = :question"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID":   &types.AttributeValueMemberS{Value: userID},
			":question": &types.AttributeValueMemberS{Value: request.Question},
		},
		ConditionExpression: aws.String("ownerID = :userID and attribute_exists(PK) and attribute_exists(SK)"),
		ReturnValues:        types.ReturnValueAllNew,
	})

	if err != nil {
		return nil, err
	}

	updatedRoom := Unmarshal(res.Attributes)

	return &updatedRoom, nil
}
