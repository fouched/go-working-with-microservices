package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"log-service/data"
	"net/http"
	"time"
)

const (
	webPort = "80"
	rpcPort = "5001"
	mongoURL = "mongodb://mongo:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	// connect to mongo
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic(err)
	}

	client = mongoClient

	// create context to disconnect
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	// close connection on app exit
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	} ()

	app := Config{
		Models: data.New(client),
	}

	// start web server
	//go app.serve()
	log.Println("Starting logging service on port: ", webPort)
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Println(err)
	}

}

//func (app *Config) serve() {
//	srv := &http.Server{
//		Addr: fmt.Sprintf(":%s", webPort),
//		Handler: app.routes(),
//	}
//
//	err := srv.ListenAndServe()
//	if err != nil {
//		log.Println(err)
//	}
//}

func connectToMongo() (*mongo.Client, error) {
	// create connection options
	clientOptions := options.Client().ApplyURI(mongoURL)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})

	// connect
	c, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Println("Error connecting:", err)
		return nil, err
	}

	log.Println("Connected to mongo")

	return c, nil
}