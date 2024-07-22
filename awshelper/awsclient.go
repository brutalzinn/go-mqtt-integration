package awshelper

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/sirupsen/logrus"
)

func clientConfig(region string) *s3.S3 {
	config := confighelper.Get()
	creds := credentials.NewStaticCredentials(config.AWS.AccessKey, config.AWS.SecretKey, "")

	sess, err := session.NewSession(&aws.Config{
		Credentials: creds,
		Region:      aws.String("us-east-1"),
		// Makes it easy to connect to localhost MinIO instances.
		S3ForcePathStyle: aws.Bool(true),
		Endpoint:         aws.String(config.AWS.EndPoint),
		LogLevel:         aws.LogLevel(aws.LogDebugWithHTTPBody),
	})
	s3Client := s3.New(sess)
	if err != nil {
		logrus.Error("Error loading AWS config: %s", err)
	}
	return s3Client
}
