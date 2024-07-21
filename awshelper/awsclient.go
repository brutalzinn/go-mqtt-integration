package awshelper

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/brutalzinn/go-mqtt-integration/confighelper"
	"github.com/sirupsen/logrus"
)

func clientConfig(region string) *s3.Client {
	configs := confighelper.Get()
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(configs.AWS.AccessKey, configs.AWS.SecretKey, "")),
		///deprecated but works for now
		config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...any) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL: configs.AWS.EndPoint,
			}, nil
		})),
	)
	if err != nil {
		logrus.Error("Error loading AWS config: %s", err)
	}

	return s3.NewFromConfig(cfg)
}
