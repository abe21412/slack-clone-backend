package util

import (
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
)

//GetSecret returns the aws secret associated with the secretName passed in
func GetSecret(secretName string) map[string]interface{} {
	region := "us-east-1"
	var data map[string]interface{}
	//Create a Secrets Manager client
	secretsManager := secretsmanager.New(session.New(), aws.NewConfig().WithRegion(region))
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := secretsManager.GetSecretValue(input)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	err = json.Unmarshal([]byte(*result.SecretString), &data)
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}
	return data
}
