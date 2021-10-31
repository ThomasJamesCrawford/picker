package environment

import (
	"context"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/mitchellh/mapstructure"
)

type Environment struct {
	CookieSecret string `mapstructure:"COOKIE_SECRET"`
}

// these will hang around for the entire life of the lambda
func New(client *ssm.Client, path *string) *Environment {
	res, err := client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           path,
		WithDecryption: true,
	})

	if err != nil {
		panic(err)
	}

	environmentMap := map[string]interface{}{}

	for _, p := range res.Parameters {
		split := strings.Split(*p.Name, "/")
		name := strings.ToUpper(split[len(split)-1])

		environmentMap[name] = p.Value
	}

	environment := &Environment{}

	decodeErr := mapstructure.Decode(&environmentMap, environment)

	if decodeErr != nil {
		panic(err)
	}

	log.Default().Println(environmentMap)

	return environment
}
