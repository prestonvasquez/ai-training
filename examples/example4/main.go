package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type index struct {
	ID   string `bson:"id"`
	Type string `bson:"type"`
}

type document struct {
	Text      string    `bson:"text"`
	Embedding []float64 `bson:"embedding"`
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// -------------------------------------------------------------------------
	// Connect to mongo

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

	db := client.Database("example4")

	// -------------------------------------------------------------------------
	// Create collection

	names, err := db.ListCollectionNames(ctx, bson.D{bson.E{Key: "name", Value: "book"}})
	if err != nil {
		return fmt.Errorf("list collections: %w", err)
	}

	if len(names) == 0 {
		fmt.Println("created book collection")

		if err := client.Database("example4").CreateCollection(ctx, "book"); err != nil {
			return fmt.Errorf("create collections: %w", err)
		}
	}

	col := client.Database("example4").Collection("book")

	// -------------------------------------------------------------------------
	// Create vector index

	indexes, err := lookupVectorIndex(ctx, col)
	if err != nil {
		return fmt.Errorf("lookupVectorIndex: %w", err)
	}

	if len(indexes) == 0 {
		if err := createVectorIndex(ctx, col.Database()); err != nil {
			return fmt.Errorf("createVectorIndex: %w", err)
		}

		fmt.Println("created vector index")

		indexes, err = lookupVectorIndex(ctx, col)
		if err != nil {
			return fmt.Errorf("lookupVectorIndex: %w", err)
		}
	}

	if len(indexes) == 0 {
		return errors.New("vector index does not exist")
	}

	fmt.Println(indexes)

	// -------------------------------------------------------------------------

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

func lookupVectorIndex(ctx context.Context, col *mongo.Collection) ([]index, error) {
	indexName := "vector_index"
	siv := col.SearchIndexes()
	cur, err := siv.List(ctx, &options.SearchIndexesOptions{Name: &indexName})
	if err != nil {
		return nil, fmt.Errorf("index: %w", err)
	}

	var indexs []index

	if err := cur.All(ctx, &indexs); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	return indexs, nil
}

func createVectorIndex(ctx context.Context, db *mongo.Database) error {
	f2 := bson.D{
		{
			Key:   "createSearchIndexes",
			Value: "book",
		},
		{
			Key: "indexes",
			Value: []bson.D{
				{
					{
						Key:   "name",
						Value: "vector_index",
					},
					{
						Key:   "type",
						Value: "vectorSearch",
					},
					{
						Key: "definition",
						Value: bson.D{
							{
								Key: "fields",
								Value: []bson.D{
									{
										{
											Key:   "type",
											Value: "vector",
										},
										{
											Key:   "numDimensions",
											Value: 300,
										},
										{
											Key:   "path",
											Value: "embedding",
										},
										{
											Key:   "similarity",
											Value: "cosine",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	res := db.RunCommand(ctx, f2)

	return res.Err()
}
