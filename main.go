package main

import (
	"context"
	"fmt"
	"github.com/Bindu483/user-api-gokit/handlers"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
)

func main() {
	fmt.Println("welcome to user-upi-gokit")
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logrus.Info("welcome to user api")

	router := mux.NewRouter()

	router.Methods(http.MethodGet).Path("/").Handler(httptransport.NewServer(handlers.WelcomeHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		msg := i.(string) //typecasting from interface to string
		_, _ = writer.Write([]byte(msg))
		return nil
	}))

	router.Methods(http.MethodPost).Path("/api/v1/users").Handler(httptransport.NewServer(handlers.CreateUserHandle(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		return nil

	}))

	router.Methods(http.MethodGet).Path("/api/v1/users").Handler(httptransport.NewServer(handlers.ListUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {

		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		_, _ = writer.Write(res)
		return nil
	}))

	router.Methods(http.MethodGet).Path("/api/v1/users/{username}").Handler(httptransport.NewServer(handlers.GetUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		_, _ = writer.Write(res)
		return nil
	}))

	router.Methods(http.MethodDelete).Path("/api/v1/users/{username}").Handler(httptransport.NewServer(handlers.DeleteUserHandler(), func(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
		return request2, nil
	}, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
		res := i.([]byte)
		_, _ = writer.Write(res)
		return nil

	}))

	server := http.Server{
		Addr:    "0.0.0.0:8500",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Error("couldnt start error")
	}
}
