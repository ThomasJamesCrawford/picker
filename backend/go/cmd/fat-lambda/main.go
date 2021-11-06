package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"picker/backend/go/pkg/environment"
	"picker/backend/go/pkg/middleware"
	"picker/backend/go/pkg/option"
	"picker/backend/go/pkg/room"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambdaV2

var client *dynamodb.Client
var ssmClient *ssm.Client

var ssmEnvironment *environment.Environment

func Handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	// This doesn't map the cookies
	// API Gateway strips the cookie header into req.Cookies, but aws-lambda-go-api-proxy doesn't seem to take this into account
	// https://github.com/awslabs/aws-lambda-go-api-proxy/issues/108
	req.Headers["cookie"] = strings.Join(req.Cookies, ",")
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
	ssmClient = ssm.NewFromConfig(cfg)

	ssmPath := os.Getenv("ssm_path")
	ssmEnvironment = environment.New(ssmClient, &ssmPath)

	r := gin.Default()

	store := cookie.NewStore([]byte(ssmEnvironment.CookieSecret))
	r.Use(sessions.Sessions(os.Getenv("session_cookie"), store))

	// Set a user ID cookie on every request
	r.Use(middleware.UserId())

	api := r.Group("/api")

	api.GET("/room/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := room.GetRoom(id, client, getUserID(c))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		if res == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if res.OwnerID != getUserID(c) {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	api.GET("/publicRoom/:id", func(c *gin.Context) {
		id := c.Param("id")
		res, err := room.GetPublicRoom(id, client, getUserID(c))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
		}

		if res == nil {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	api.GET("/publicRoom/:id/available", func(c *gin.Context) {
		id := c.Param("id")
		res, err := room.GetPublicRoom(id, client, getUserID(c))

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		if res == nil {
			c.JSON(http.StatusOK, gin.H{"available": true})
			return
		}

		c.JSON(http.StatusOK, gin.H{"available": false})
	})

	api.POST("/room", func(c *gin.Context) {
		createRoomRequest := &room.CreateRoomRequest{}

		err = c.ShouldBindJSON(&createRoomRequest)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		room, err := room.NewRoom(createRoomRequest, getUserID(c), client)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, room)
	})

	api.PATCH("/room/:roomID/option/:optionID/select", func(c *gin.Context) {
		roomID := c.Param("roomID")
		optionID := c.Param("optionID")

		res, err := option.SelectOption(optionID, getUserID(c), roomID, client)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	api.PATCH("/room/:roomID/option/:optionID/unselect", func(c *gin.Context) {
		roomID := c.Param("roomID")
		optionID := c.Param("optionID")

		res, err := option.UnselectOption(optionID, getUserID(c), roomID, client)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	api.DELETE("/room/:roomID/option/:optionID", func(c *gin.Context) {
		roomID := c.Param("roomID")
		optionID := c.Param("optionID")

		res, err := option.Delete(optionID, getUserID(c), roomID, client)

		if err != nil {
			c.AbortWithError(http.StatusForbidden, err)
			return
		}

		c.JSON(http.StatusOK, res)
	})

	ginLambda = ginadapter.NewV2(r)
}

func getUserID(c *gin.Context) string {
	return fmt.Sprintf("%v", sessions.Default(c).Get("user_id"))
}

func main() {
	lambda.Start(Handler)
}
