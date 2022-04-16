package appconfig

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/joho/godotenv"
)

type Config struct {
	PortApp          string
	AwsConfig        *aws.Config
	isDevelopment    bool
	localstackRegion string
	awsPartitionId   string
	localstackUrl    string
}

func LoadConfig() Config {

	godotenv.Load()
	configApp := Config{}
	env := os.Getenv("ENVIRONMENT")
	configApp.PortApp = os.Getenv("PORT")
	configApp.localstackRegion = os.Getenv("LOCALSTACK_AWS_REGION")
	configApp.awsPartitionId = os.Getenv("LOCALSTACK_PARTITION_ID")
	configApp.localstackUrl = os.Getenv("LOCALSTACK_URL")
	configApp.PortApp = os.Getenv("PORT")

	if env == "dev" {

		configApp.isDevelopment = true
		if cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   configApp.awsPartitionId,
					URL:           configApp.localstackUrl,
					SigningRegion: configApp.AwsConfig.Region,
				}, nil
			}))); err != nil {
			panic(err)
		} else {
			configApp.AwsConfig = &cfg
		}

	} else {

		configApp.isDevelopment = false
		if cfg, err := config.LoadDefaultConfig(context.Background()); err != nil {
			panic(err)
		} else {
			configApp.AwsConfig = &cfg
		}
	}

	return configApp
}
