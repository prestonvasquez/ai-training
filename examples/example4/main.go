package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/vector/foundation/mongodb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// -------------------------------------------------------------------------
	// Connect to mongo

	client, err := mongodb.Connect(ctx, "mongodb://localhost:27017", "ardan", "ardan")
	if err != nil {
		return fmt.Errorf("connectToMongo: %w", err)
	}
	defer client.Disconnect(ctx)

	fmt.Println("Connected to MongoDB")

	// -------------------------------------------------------------------------
	// Create database and collection

	const dbName = "example4"
	const collectionName = "book"

	db := client.Database(dbName)

	col, err := mongodb.CreateCollection(ctx, db, collectionName)
	if err != nil {
		return fmt.Errorf("createCollection: %w", err)
	}

	fmt.Println("Created Collection")

	// -------------------------------------------------------------------------
	// Create vector index

	const indexName = "vector_index"

	settings := mongodb.VectorIndexSettings{
		NumDimensions: 4,
		Path:          "embedding",
		Similarity:    "cosine",
	}

	if err := mongodb.CreateVectorIndex(ctx, col, indexName, settings); err != nil {
		return fmt.Errorf("createVectorIndex: %w", err)
	}

	fmt.Println("Created Vector Index")

	// -------------------------------------------------------------------------
	// Store some documents with their embeddings.

	if err := storeDocuments(ctx, col); err != nil {
		return fmt.Errorf("storeDocuments: %w", err)
	}

	if err := queryDocuments(ctx, col); err != nil {
		return fmt.Errorf("storeDocuments: %w", err)
	}

	return nil
}

func storeDocuments(ctx context.Context, col *mongo.Collection) error {
	col.DeleteMany(ctx, bson.D{})

	type document struct {
		ID        int       `bson:"id"`
		Text      string    `bson:"text"`
		Embedding []float64 `bson:"embedding"`
	}

	d1 := document{
		ID:        1,
		Text:      "this is text 1",
		Embedding: []float64{1.0, 2.0, 3.0, 4.0},
	}

	res, err := col.InsertOne(ctx, d1)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	fmt.Println(res.InsertedID)

	d2 := document{
		ID:        2,
		Text:      "this is text 2",
		Embedding: []float64{1.5, 2.5, 3.5, 4.5},
	}

	res, err = col.InsertOne(ctx, d2)
	if err != nil {
		return fmt.Errorf("insert: %w", err)
	}

	fmt.Println(res.InsertedID)

	return nil
}

func queryDocuments(ctx context.Context, col *mongo.Collection) error {
	findRes, err := col.Find(ctx, bson.D{})
	if err != nil {
		return fmt.Errorf("find: %w", err)
	}

	var r []struct {
		Text      string    `bson:"text"`
		Embedding []float64 `bson:"embedding"`
		Score     float64   `bson:"score"`
	}

	fmt.Println("BATCH", findRes.RemainingBatchLength())

	if err := findRes.All(ctx, &r); err != nil {
		return fmt.Errorf("find: %w", err)
	}
	findRes.Close(ctx)

	fmt.Println("find:", r)

	fmt.Println("---- VECTOR QUERY ----")

	// db.book.aggregate([ { "$vectorSearch": { "index": "vector_index", "exact": false, "numCandidates": 10, "path": "embedding", "queryVector": [1.2, 2.2, 3.2, 4.2], "limit": 10 } }, { "$project": { "text": 1, "embedding": 1, "score": { "$meta": "vectorSearchScore" } } }])

	pipeline := mongo.Pipeline{
		{{
			Key: "$vectorSearch",
			Value: bson.M{
				"index":         "vector_index",
				"exact":         false,
				"path":          "embedding",
				"queryVector":   []float64{1.2, 2.2, 3.2, 4.2},
				"numCandidates": 10,
				"limit":         10,
			}},
		},
		{{
			Key: "$project",
			Value: bson.M{
				"text":      1,
				"embedding": 1,
				"score": bson.M{
					"$meta": "vectorSearchScore",
				},
			}},
		},
	}

	cur, err := col.Aggregate(ctx, pipeline)
	if err != nil {
		return fmt.Errorf("aggregate: %w", err)
	}

	fmt.Println("BATCH", cur.RemainingBatchLength())

	var results []struct {
		Text      string    `bson:"text"`
		Embedding []float64 `bson:"embedding"`
		Score     float64   `bson:"score"`
	}
	if err := cur.All(ctx, &results); err != nil {
		return fmt.Errorf("all: %w", err)
	}
	cur.Close(ctx)

	fmt.Println(results)

	return nil
}
