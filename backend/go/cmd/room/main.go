package main

import (
	"context"
	"picker/backend/go/pkg/room"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	return ginLambda.ProxyWithContext(ctx, req)
}

func init() {
	r := gin.Default()

	r.GET("/room/:id", func(c *gin.Context) {
		c.JSON(200, &room.Room{})
	})

	ginLambda = ginadapter.NewV2(r)
}

func main() {
	lambda.Start(Handler)
}
