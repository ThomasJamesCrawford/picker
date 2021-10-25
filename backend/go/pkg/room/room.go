package room

type Room struct {
	PK string `dynamodbav:"PK"`
	SK string `dynamodbav:"SK"`

	Type string `dynamodbav:"type"`

	Name       string `json:"name"`
	SecretName string `json:"secretName"`
}
