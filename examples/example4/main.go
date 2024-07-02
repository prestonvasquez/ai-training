package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type document struct {
	Text             string
	VectorEmbeddings []float64
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	auth := options.Client().SetAuth(options.Credential{
		Username: "ardan",
		Password: "ardan",
	})

	uri := options.Client().ApplyURI("mongodb://localhost:27017")

	client, err := mongo.Connect(ctx, auth, uri)
	if err != nil {
		return fmt.Errorf("connect: %w", err)
	}
	defer client.Disconnect(ctx)

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("ping: %w", err)
	}

	fmt.Println("Connected to MongoDB")

	col := client.Database("example4").Collection("book")

	d1 := document{
		Text:             "this is text 1",
		VectorEmbeddings: []float64{1.0, 2.0, 3.0, 4.0},
	}

	ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	res, err := col.InsertOne(ctx1, d1)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	fmt.Println(res)

	return nil
}
