// This example show you how to use MongoDB and Ollama to perform a vector
// search for a user question. The search will return the top 5 chunks from
// the database. Then these chunks are sent to the Llama model to create a
// coherent response.
//
// # Running the example:
//
//	$ make example6
//
// # This requires running the following commands:
//
//	$ make dev-up   // This starts the mongodb and ollama service in docker compose.
//	$ make example5 // This creates the book.embeddings file

package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ardanlabs/ai-training/foundation/mongodb"
	"github.com/tmc/langchaingo/llms/ollama"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type searchResult struct {
	ID        int       `bson:"id"`
	Text      string    `bson:"text"`
	Embedding []float64 `bson:"embedding"`
	Score     float64   `bson:"score"`
}

// =============================================================================

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Open a connection with ollama to access the model.
	llm, err := ollama.New(ollama.WithModel("mxbai-embed-large"))
	if err != nil {
		return fmt.Errorf("ollama: %w", err)
	}

	// Connect to mongodb.
	client, err := mongodb.Connect(ctx, "mongodb://localhost:27017", "ardan", "ardan")
	if err != nil {
		return fmt.Errorf("connectToMongo: %w", err)
	}

	const dbName = "example5"
	const collectionName = "book"

	// Capture a connection to the collection. We assume this exists with
	// data already.
	col := client.Database(dbName).Collection(collectionName)

	// -------------------------------------------------------------------------

	question := "what is an interface?"

	results, err := vectorSearch(ctx, llm, col, question)
	if err != nil {
		return fmt.Errorf("vectorSearch: %w", err)
	}

	for _, res := range results {
		fmt.Printf("ID: %d, Score: %.3f%%\n%s\n\n", res.ID, res.Score*100, res.Text)
	}

	return nil
}

func vectorSearch(ctx context.Context, llm *ollama.LLM, col *mongo.Collection, question string) ([]searchResult, error) {

	// Get the vector embedding for the question.
	embedding, err := llm.CreateEmbedding(context.Background(), []string{question})
	if err != nil {
		return nil, fmt.Errorf("create embedding: %w", err)
	}

	// We want to find the nearest neighbors from the specified vector.
	pipeline := mongo.Pipeline{
		{{
			Key: "$vectorSearch",
			Value: bson.M{
				"index":         "vector_index",
				"exact":         false,
				"path":          "embedding",
				"queryVector":   embedding[0],
				"numCandidates": 5,
				"limit":         5,
			}},
		},
		{{
			Key: "$project",
			Value: bson.M{
				"id":        1,
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
		return nil, fmt.Errorf("aggregate: %w", err)
	}
	defer cur.Close(ctx)

	var results []searchResult
	if err := cur.All(ctx, &results); err != nil {
		return nil, fmt.Errorf("all: %w", err)
	}

	return results, nil
}
