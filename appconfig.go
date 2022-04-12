package appconfig

import (
	"context"
	"log"
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

	configApp := Config{}
	godotenv.Load()

	env := os.Getenv("ENVIRONMENT")
	if env == "dev" {

		if err := godotenv.Load(".env"); err != nil {
			log.Fatalf("Error loading .env file for development environment")
		}

		configApp.isDevelopment = true
		configApp.PortApp = os.Getenv("PORT")
		configApp.localstackRegion = os.Getenv("LOCALSTACK_AWS_REGION")
		configApp.awsPartitionId = os.Getenv("LOCALSTACK_PARTITION_ID")
		configApp.localstackUrl = os.Getenv("LOCALSTACK_URL")

		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					PartitionID:   configApp.awsPartitionId,
					URL:           configApp.localstackUrl,
					SigningRegion: configApp.AwsConfig.Region,
				}, nil
			})))

		if err != nil {
			panic(err)
		}

		configApp.AwsConfig = &cfg

	} else {

		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			panic(err)
		}

		port := os.Getenv("PORT")
		if port == "" {
			panic("the port number cannot be empty")
		}

		configApp.PortApp = port
		configApp.AwsConfig = &cfg
		configApp.AwsConfig = &aws.Config{}
		configApp.isDevelopment = false

	}

	return configApp
}
