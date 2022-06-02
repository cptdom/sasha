package utils

import (
	"log"

	"github.com/aws/aws-sdk-go/aws/credentials"
)

func FetchCredsFromEnv() interface{} {
	creds := credentials.NewEnvCredentials()
	credVal, err := creds.Get()
	if err != nil {
		log.Fatalf("No creds found, exiting: %v", err)
	}
	return credVal
}
