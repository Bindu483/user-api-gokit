package handlers

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
	"os"
)

func WelcomeHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		hostName, err := os.Hostname()
		if err != nil {
			logrus.WithError(err).Error("couldnot fetch hostname")
		}
		return "welcome to user-api-gokit " + hostName, err
	}
}
