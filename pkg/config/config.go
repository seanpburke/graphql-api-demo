package config

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/spf13/viper"
)

const (
	defaultConfigPath    = "."
	defaultListenAddress = ":8080"
	envConfigDirectory   = "CONFIG_DIR"
	envDDBEndpoint       = "AWS_DDB_ENDPOINT"
)

var configPath = defaultConfigPath

var Config struct {
	// These fields come from config.json
	AppName  string // This app's name
	Region   string // AWS Region
	Registry string // AWS ECR registry
	Table    string // AWS DynamoDB Table Name
	Listen   string // HTTP Listener Address

	AWS *session.Session
	DDB *dynamodb.DynamoDB
}

func init() {

	viper.AutomaticEnv()
	viper.SetConfigName("config")
	viper.SetDefault("Listen", defaultListenAddress)

	if v := os.Getenv(envConfigDirectory); v != "" {
		viper.AddConfigPath(v)
	} else {
		viper.AddConfigPath(".")
	}
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("%s", err)
	}
	if err := viper.Unmarshal(&Config); err != nil {
		log.Fatalf("%s", err)
	}

	// Initialize the AWS session. In order to load credentials
	// from the shared credentials file ~/.aws/credentials,
	// and region from the shared configuration file ~/.aws/config,
	// set AWS_SDK_LOAD_CONFIG=1.
	cnf := aws.NewConfig().
		WithRegion(Config.Region).
		WithCredentialsChainVerboseErrors(true)
	Config.AWS = session.Must(session.NewSession(cnf))

	// Create DynamoDB client
	// The service Endpoint can be specified in an environment variable.
	// This enables you to operate on the local DynamoDB, but obviously
	// this is of no use for Lambda functions.
	if endpoint := os.Getenv(envDDBEndpoint); endpoint != "" {
		cnf = cnf.WithEndpoint(endpoint)
	}
	Config.DDB = dynamodb.New(Config.AWS, cnf)
}
