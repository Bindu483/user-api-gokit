package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

func CreateUserHandle() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		cli, err := clientV3.New(clientV3.Config{
			Endpoints:   []string{fmt.Sprintf("%s:2379", etcdIp)},
			DialTimeout: 5 * time.Second,
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
		})
		if err != nil {
			// handle error!
			logrus.WithError(err).Error("couldn't connect to database")
			return nil, err

		}
		defer cli.Close()
		req := request.(*http.Request) //typecasting from interface to http request type
		var userObject User
		err = json.NewDecoder(req.Body).Decode(&userObject)
		if err != nil {
			logrus.WithError(err).Error("couldn't decode user object")
			return nil, err

		}
		if userObject.FirstName == "" {
			return nil, errors.New("first name is empty")

		}
		userByte, err := json.Marshal(userObject)
		if err != nil {
			logrus.WithError(err).Error("count marshal user object")
			return nil, err
		}
		// Key => /name/bindu
		// value => json
		_, err = cli.Put(context.Background(), fmt.Sprintf("/name/%s", userObject.FirstName), string(userByte))
		if err != nil {
			logrus.WithError(err).Error("couldnt write to the database")
			return nil, err
		}
		return "user created successfully", nil
	}
}
