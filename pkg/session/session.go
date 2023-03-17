package session

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	logger "github.com/sirupsen/logrus"
)

var (
	sessionConfig           *aws.Config
	sessionConfigWithRegion *aws.Config
)

func GetSessionConfig(ctx context.Context) (aws.Config, bool) {
	if sessionConfig != nil {
		return *sessionConfig, true
	}

	cfg, success := getNewSession(ctx, "ap-southeast-1") // SG
	if !success {
		return aws.Config{}, false
	}
	sessionConfig = &cfg

	return *sessionConfig, true
}

func GetSessionConfigWithRegion(ctx context.Context, region string) (aws.Config, bool) {
	if sessionConfigWithRegion != nil {
		if sessionConfigWithRegion.Region == region {
			return *sessionConfigWithRegion, true
		}
	}

	cfg, success := getNewSession(ctx, region)
	if !success {
		return aws.Config{}, false
	}
	sessionConfigWithRegion = &cfg
	return *sessionConfigWithRegion, true
}

func getNewSession(ctx context.Context, region string) (aws.Config, bool) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		logger.WithFields(logger.Fields{
			"error": err.Error(),
			"code":  "CFGErr",
		}).Error("failed to load aws config")
		return aws.Config{}, false
	}

	return cfg, true
}
