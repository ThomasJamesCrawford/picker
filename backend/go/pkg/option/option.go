package option

import (
	"context"
	"fmt"
	"os"
	"picker/backend/go/pkg/dynamodbTypes"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/twinj/uuid"
)

type Option struct {
	// DynamoDB
	PK     string `dynamodbav:"PK" json:"-"`
	SK     string `dynamodbav:"SK" json:"-"`
	GSI1PK string `dynamodbav:"GSI1PK" json:"-"`
	GSI1SK string `dynamodbav:"GSI1SK" json:"-"`
	Type   string `dynamodbav:"type" json:"-"`

	// Public
	ID     string `json:"id"`
	RoomID string `json:"roomID"`

	// Private
	SelectedByID *string `dynamodbav:"selectedByID,omitEmpty" json:"-"`
}

func NewOption(option string, userID string, roomID string, client *dynamodb.Client) *Option {
	optionID := uuid.NewV4().String()

	return &Option{
		PK: fmt.Sprintf("ROOM#%s", roomID),
		// This lets us use a BEGINS_WITH in our single table to pull in a room and all the options with one query
		SK:     fmt.Sprintf("ROOM_OPTION#%s", optionID),
		GSI1PK: fmt.Sprintf("USER#%s", userID),
		GSI1SK: fmt.Sprintf("ROOM_OPTION#%s", optionID),
		Type:   dynamodbTypes.Option,

		ID:     optionID,
		RoomID: roomID,

		SelectedByID: nil,
	}
}

func BatchWriteOptions(options []*Option, client *dynamodb.Client) error {
	// BatchWriteItem does max 25
	chunkedOptions := chunk(options, 25)

	// TODO loop could be async
	for _, chunk := range chunkedOptions {
		var items []types.WriteRequest

		for _, option := range chunk {
			item, err := attributevalue.MarshalMap(option)

			if err != nil {
				panic(err)
			}

			items = append(items, types.WriteRequest{
				PutRequest: &types.PutRequest{
					Item: item,
				},
			})
		}

		request := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]types.WriteRequest{
				os.Getenv("table"): items,
			},
		}

		_, err := client.BatchWriteItem(context.TODO(), request)

		if err != nil {
			return err
		}
	}

	return nil
}

// https://freshman.tech/snippets/go/split-slice-into-chunks/#loop-through-the-number-of-chunks
func chunk(options []*Option, chunkSize int) [][]*Option {
	var chunks [][]*Option

	for i := 0; i < len(options); i += chunkSize {
		end := i + chunkSize

		if end > len(options) {
			end = len(options)
		}

		chunks = append(chunks, options[i:end])
	}

	return chunks
}

func Unmarshal(item map[string]types.AttributeValue) *Option {
	option := &Option{}
	if err := attributevalue.Unmarshal(item, option); err != nil {
		panic(err)
	}

	return option
}
