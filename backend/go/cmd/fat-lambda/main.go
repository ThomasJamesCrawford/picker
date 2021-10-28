package main

import (
	"context"
	"net/http"
	"os"
	"picker/backend/go/pkg/room"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2

var client *dynamodb.Client

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	cfg, err := config.LoadDefaultConfig(context.TODO(), func(o *config.LoadOptions) error {
		o.Region = os.Getenv("region")
		return nil
	})

	if err != nil {
		panic(err)
	}

	client = dynamodb.NewFromConfig(cfg)

	r := gin.Default()

	r.GET("/publicRoom/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := room.GetPublicRoom(id, client)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
		}

		if res == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	r.GET("/publicRoom/:id/available", func(c *gin.Context) {
		id := c.Param("id")
		res, err := room.GetPublicRoom(id, client)

		if err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if res == nil {
			c.JSON(http.StatusOK, gin.H{"available": true})
			return
		}

		c.JSON(http.StatusOK, gin.H{"available": false})
	})

	ginLambda = ginadapter.NewV2(r)
}

func main() {
	lambda.Start(Handler)
}
