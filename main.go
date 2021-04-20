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

	router.Methods(http.MethodGet).Path("/").
		Handler(httptransport.
			NewServer(handlers.WelcomeHandler(),
				defaultRequestDecoder,
				func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
					msg := i.(string) //typecasting from interface to string
					_, _ = writer.Write([]byte(msg))
					return nil
				}))

	router.Methods(http.MethodPost).
		Path("/api/v1/users").
		Handler(httptransport.
			NewServer(handlers.CreateUserHandle(),
				defaultRequestDecoder, func(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
					return nil
				}))

	router.Methods(http.MethodGet).
		Path("/api/v1/users").
		Handler(httptransport.NewServer(handlers.ListUserHandler(),
			defaultRequestDecoder,
			encodeByteSlice))

	router.Methods(http.MethodGet).
		Path("/api/v1/users/{username}").
		Handler(httptransport.NewServer(handlers.GetUserHandler(),
			defaultRequestDecoder,
			encodeByteSlice))

	router.Methods(http.MethodDelete).Path("/api/v1/users/{username}").
		Handler(httptransport.
			NewServer(handlers.DeleteUserHandler(),
				defaultRequestDecoder,
				encodeByteSlice))

	server := http.Server{
		Addr:    "0.0.0.0:8500",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		logrus.WithError(err).Error("couldnt start error")
	}
}

func defaultRequestDecoder(ctx context.Context, request2 *http.Request) (request interface{}, err error) {
	return request2, nil
}

func encodeByteSlice(ctx context.Context, writer http.ResponseWriter, i interface{}) error {
	res := i.([]byte)
	_, _ = writer.Write(res)
	return nil
}
