package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)
import "github.com/sirupsen/logrus"
import "github.com/gorilla/mux"
import httptransport "github.com/go-kit/kit/transport/http"
import "go.etcd.io/etcd/client/v3"

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

func main() {
	fmt.Println("welcome to user-upi-gokit")
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	logrus.Info("welcome to user api")
	router := mux.NewRouter()
	router.Methods(http.MethodGet).Path("/").Handler(httptransport.NewServer(welcomeHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		msg := i.(string) //typecasting from interface to string
		writer.Write([]byte(msg))
		return nil
	}))
	router.Methods(http.MethodPost).Path("/api/v1/users").Handler(httptransport.NewServer(createUserHandle(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		return nil

	}))
	router.Methods(http.MethodGet).Path("/api/v1/users").Handler(httptransport.NewServer(listUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {

		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		writer.Write(res)
		return nil
	}))
	router.Methods(http.MethodGet).Path("/api/v1/users/{username}").Handler(httptransport.NewServer(getUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		writer.Write(res)
		return nil

	}))
	router.Methods(http.MethodDelete).Path("/api/v1/users/{username}").Handler(httptransport.NewServer(deleteUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		writer.Write(res)
		return nil

	}))
	server := http.Server{
		Addr: "0.0.0.0:8500", Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Error("couldnt start error")
	}
}

func deleteUserHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		cli, err := clientv3.New(clientv3.Config{
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
		deleteResponse, err := cli.Delete(context.Background(), fmt.Sprintf("/name/%s", userName))
		if err != nil {
			logrus.WithError(err).Error("couldnt read username")
			return nil, err
		}
		if deleteResponse.Deleted != 0 {
			return []byte("no user deleted"), nil
		}

		return []byte(" user deleted"), nil

	}
}

func getUserHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		cli, err := clientv3.New(clientv3.Config{
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
			userObject := user{}
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

func listUserHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		logrus.Info(etcdIp)
		if etcdIp == "" {
			logrus.Error("etcdip was empty")
			os.Exit(1)
		}
		cli, err := clientv3.New(clientv3.Config{
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
		getResponse, err := cli.Get(context.Background(), "/name/", clientv3.WithPrefix())
		if err != nil {
			logrus.WithError(err).Error("couldnt fetch users")
			return nil, err
		}
		var userList []user
		for _, kv := range getResponse.Kvs {
			userObject := user{}
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

func createUserHandle() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		etcdIp := os.Getenv("ETCD_IP")
		cli, err := clientv3.New(clientv3.Config{
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
		req := request.(*http.Request) //typecasting from interface to http request type
		var userObject user
		err = json.NewDecoder(req.Body).Decode(&userObject)
		if err != nil {
			logrus.WithError(err).Error("couldnot decode userobject")
			return nil, err

		}
		if userObject.FirstName == "" {
			return nil, errors.New("first name is empty")

		}
		userByte, err := json.Marshal(userObject)
		if err != nil {
			logrus.WithError(err).Error("coudnt marshal userobject")
			return nil, err
		}
		// Key => /name/bindu
		// value => json
		_, err = cli.Put(context.Background(), fmt.Sprintf("/name/%s", userObject.FirstName), string(userByte))
		if err != nil {
			logrus.WithError(err).Error("couldnt write to the database")
			return nil, err
		}
		return "user created succesfully", nil
	}
}

func welcomeHandler() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		hostName, err := os.Hostname()
		if err != nil {
			logrus.WithError(err).Error("couldnot fetch hostname")
		}
		return "welcome to user-api-gokit " + hostName, err
	}

}
