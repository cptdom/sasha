package s3

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetSession(endpoint string) *session.Session {
	sess := session.Must(session.NewSession(&aws.Config{
		Region:   aws.String("default"),
		Endpoint: aws.String(endpoint),
	}))
	return sess
}

type SessionClient struct {
	Session session.Session
	Client  s3.S3
}

func (s *SessionClient) GetAccessID() string {
	sessionData, err := s.Session.Config.Credentials.Get()
	if err != nil {
		fmt.Println("No AccessID found.")
		os.Exit(1)
	}
	return sessionData.AccessKeyID
}

func (s *SessionClient) GetEndpoint() string {
	endpoint := *s.Session.Config.Endpoint
	if strings.HasPrefix(endpoint, "https://") {
		return strings.Split(endpoint, "//")[1]
	}
	return endpoint
}

func (s *SessionClient) Whoami() {
	fmt.Println(s.GetAccessID())
}
