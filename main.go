package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-kit/kit/endpoint"
	"net/http"
)
import "github.com/sirupsen/logrus"
import "github.com/gorilla/mux"
import httptransport "github.com/go-kit/kit/transport/http"

type user struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
}

var DB []*user

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
	server := http.Server{
		Addr: "0.0.0.0:8500", Handler: router,
	}
	err := server.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Error("couldnt start error")
	}
}

func createUserHandle() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		var userObject user
		err := json.NewDecoder(request.Body).Decode(&userObject)
		if err != nil {
			fmt.Println("coudnt  read requestBody", err)
			writer.Write([]byte("couldnt read requestBody"))
			return

		}
		DB = append(DB, &userObject)
	}
}
