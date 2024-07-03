package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/vector/foundation/mongodb"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// -------------------------------------------------------------------------
	// Connect to mongo

	client, err := mongodb.Connect(ctx, "mongodb://localhost:27017")
	if err != nil {
		return fmt.Errorf("connectToMongo: %w", err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB")

	// -------------------------------------------------------------------------
	// Create collection

	const dbName = "example4"
	const collectionName = "book"

	col, err := mongodb.CreateCollection(ctx, client, dbName, collectionName)
	if err != nil {
		return fmt.Errorf("createCollection: %w", err)
	}

	fmt.Println("Created Collection")

	// -------------------------------------------------------------------------
	// Create vector index

	settings := mongodb.VectorIndexSettings{
		NumDimensions: 300,
		Path:          "embedding",
		Similarity:    "cosine",
	}

	if err := mongodb.CreateVectorIndex(ctx, col, "vector_index", settings); err != nil {
		return fmt.Errorf("createVectorIndex: %w", err)
	}

	fmt.Println("Created Vector Index")

	// -------------------------------------------------------------------------
	// Store some documents with their embeddings.

	// type document struct {
	// 	Text      string    `bson:"text"`
	// 	Embedding []float64 `bson:"embedding"`
	// }

	// fmt.Println(res)

	// d1 := document{
	// 	Text:      "this is text 1",
	// 	Embedding: []float64{1.0, 2.0, 3.0, 4.0},
	// }

	// ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()

	// res, err := col.InsertOne(ctx1, d1)
	// if err != nil {
	// 	return fmt.Errorf("insert: %w", err)
	// }

	//col.FindOne(ctx, filter interface{}, opts ...*options.FindOneOptions)

	return nil
}
