package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

func GetUserHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		cli, err := clientV3.New(clientV3.Config{
			Endpoints:   []string{fmt.Sprintf("%s:2379", etcdIp)},
			DialTimeout: 5 * time.Second,
			DialOptions: []grpc.DialOption{grpc.WithBlock()},
		})
		if err != nil {
			// handle error!
			logrus.WithError(err).Error("coudnt connect to database")
			return nil, err

		}
		defer cli.Close()
		req := request.(*http.Request)
		urlVars := mux.Vars(req)
		userName := urlVars["username"]
		getResponse, err := cli.Get(context.Background(), fmt.Sprintf("/name/%s", userName))
		if err != nil {
			logrus.WithError(err).Error("couldnt read username")
			return nil, err
		}
		if len(getResponse.Kvs) > 0 {
			userDbResponse := getResponse.Kvs[0]
			userObject := User{}
			err = json.Unmarshal(userDbResponse.Value, &userObject)
			if err != nil {
				logrus.WithError(err).Error("couldnt unmarshal userobject")
				return nil, err
			}
			responseByte, err := json.Marshal(userObject)
			if err != nil {
				logrus.WithError(err).Error("couldnt marshal userobject")
				return nil, err
			}
			return responseByte, nil

		}
		return []byte("no user found"), nil

	}
}
