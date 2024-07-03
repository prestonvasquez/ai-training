package mongodb

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Connect attempts to connect to a mongo db instance.
func Connect(ctx context.Context, host string, userName string, password string) (*mongo.Client, error) {
	auth := options.Client().SetAuth(options.Credential{
		Username: userName,
		Password: password,
	})

	uri := options.Client().ApplyURI(host)

	client, err := mongo.Connect(ctx, auth, uri)
	if err != nil {
		return nil, fmt.Errorf("connect: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}

	return client, nil
}

// CreateCollection will create the specified collection in the specified
// database if it doesn't already exist.
func CreateCollection(ctx context.Context, db *mongo.Database, collectionName string) (*mongo.Collection, error) {
	names, err := db.ListCollectionNames(ctx, bson.D{bson.E{Key: "name", Value: collectionName}})
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}

	if len(names) == 0 {
		if err := db.CreateCollection(ctx, collectionName); err != nil {
			return nil, fmt.Errorf("create collections: %w", err)
		}
	}

	return db.Collection(collectionName), nil
}

// CreateVectorIndex creates a very specific vector index for our example.
func CreateVectorIndex(ctx context.Context, col *mongo.Collection, vectorIndexName string, settings VectorIndexSettings) error {
	indexes, err := lookupVectorIndex(ctx, col)
	if err != nil {
		return fmt.Errorf("lookupVectorIndex: %w", err)
	}

	if len(indexes) == 0 {
		if err := runCreateIndexCmd(ctx, col, vectorIndexName, settings); err != nil {
			return fmt.Errorf("createVectorIndex: %w", err)
		}

		indexes, err = lookupVectorIndex(ctx, col)
		if err != nil {
			return fmt.Errorf("lookupVectorIndex: %w", err)
		}
	}

	if len(indexes) == 0 {
		return errors.New("vector index does not exist")
	}

	return nil
}

// =============================================================================

func lookupVectorIndex(ctx context.Context, col *mongo.Collection) ([]Index, error) {
	indexName := "vector_index"
	siv := col.SearchIndexes()
	cur, err := siv.List(ctx, &options.SearchIndexesOptions{Name: &indexName})
	if err != nil {
		return nil, fmt.Errorf("index: %w", err)
	}

	var indexs []Index

	if err := cur.All(ctx, &indexs); err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	return indexs, nil
}

func runCreateIndexCmd(ctx context.Context, col *mongo.Collection, vectorIndexName string, settings VectorIndexSettings) error {
	fields := bson.D{
		{
			Key:   "type",
			Value: "vector",
		},
		{
			Key:   "numDimensions",
			Value: settings.NumDimensions,
		},
		{
			Key:   "path",
			Value: settings.Path,
		},
		{
			Key:   "similarity",
			Value: settings.Similarity,
		},
	}

	idx := bson.D{
		{
			Key:   "createSearchIndexes",
			Value: col.Name(),
		},
		{
			Key: "indexes",
			Value: []bson.D{
				{
					{
						Key:   "name",
						Value: vectorIndexName,
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
									fields,
								},
							},
						},
					},
				},
			},
		},
	}

	res := col.Database().RunCommand(ctx, idx)

	return res.Err()
}
