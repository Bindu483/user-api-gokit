package handlers

import (
	"context"
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

func DeleteUserHandler() endpoint.Endpoint {
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
		req := request.(*http.Request)
		urlVars := mux.Vars(req)
		userName := urlVars["username"]
		deleteResponse, err := cli.Delete(context.Background(), fmt.Sprintf("/name/%s", userName))
		if err != nil {
			logrus.WithError(err).Error("couldn't read username")
			return nil, err
		}
		if deleteResponse.Deleted != 0 {
			return []byte("no user deleted"), nil
		}

		return []byte(" user deleted"), nil

	}
}
