package main

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

func welcomeHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		return "welcome to user-api-gokit", err
	}

}
