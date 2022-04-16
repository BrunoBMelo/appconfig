package appconfig

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Config struct {
	PortApp       string
	AwsConfig     *aws.Config
	isDevelopment bool
}

func LoadConfig() Config {

	configApp := Config{}
	env := os.Getenv("ENVIRONMENT")
	configApp.PortApp = os.Getenv("PORT")

	if env == "dev" {

		fmt.Println("Loading variables to environment: DEV")
		configApp.isDevelopment = true
		if cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   os.Getenv("LOCALSTACK_PARTITION_ID"),
					URL:           os.Getenv("LOCALSTACK_URL"),
					SigningRegion: os.Getenv("LOCALSTACK_AWS_REGION"),
				}, nil
			}))); err != nil {
			panic(err)
		} else {
			configApp.AwsConfig = &cfg
		}

	} else {
		fmt.Println("Loading variables to environment: PROD")
		configApp.isDevelopment = false
		if cfg, err := config.LoadDefaultConfig(context.Background()); err != nil {
			panic(err)
		} else {
			configApp.AwsConfig = &cfg
		}
	}

	return configApp
}
