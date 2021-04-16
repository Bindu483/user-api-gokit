package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"github.com/sirupsen/logrus"
	clientV3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"os"
	"time"
)

func ListUserHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		logrus.Info(etcdIp)
		if etcdIp == "" {
			logrus.Error("etcdip was empty")
			os.Exit(1)
		}
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
		logrus.Info("client connection succesfull", cli.ActiveConnection().GetState().String())

		defer cli.Close()
		/*
			/name/bindu
			/name/something
			/name/somethingwlse
			/cars/huyndai

			/name --prefix
		*/
		getResponse, err := cli.Get(context.Background(), "/name/", clientV3.WithPrefix())
		if err != nil {
			logrus.WithError(err).Error("couldnt fetch users")
			return nil, err
		}
		var userList []User
		for _, kv := range getResponse.Kvs {
			userObject := User{}
			err := json.Unmarshal(kv.Value, &userObject)
			if err != nil {
				logrus.WithError(err).Error("coudnt unmarshal userObject")
				return nil, err
			}
			userList = append(userList, userObject)

		}
		byteArray, err := json.Marshal(userList)
		if err != nil {
			logrus.WithError(err).Error("couldnt fetch users")
			return nil, err
		}
		logrus.Info(string(byteArray))
		return byteArray, nil
	}
}
