package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ardanlabs/vector/foundation/vector"
	"github.com/tmc/langchaingo/llms/ollama"
)

/*
	https://machinelearningmastery.com/what-are-word-embeddings/
	https://machinelearningmastery.com/use-word-embedding-layers-deep-learning-keras/

	The position of a data point in the learned vector space is referred to as
	its	embedding.
*/

type data struct {
	Name string
	Text string
	Vec  []float32
}

// Vector can convert the specified data into a vector.
func (d data) Vector() []float32 {
	return d.Vec
}

// =============================================================================

func main() {
	llm, err := ollama.New(ollama.WithModel("mxbai-embed-large"))
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------------

	// Apply the feature vectors to the hand crafted dataPoints.
	dataPoints := []vector.Data{
		data{Name: "Horse   ", Text: "Animal, Female"},
		data{Name: "Man     ", Text: "Human,  Male,   Pants, Poor, Worker"},
		data{Name: "Woman   ", Text: "Human,  Female, Dress, Poor, Worker"},
		data{Name: "King    ", Text: "Human,  Male,   Pants, Rich, Ruler"},
		data{Name: "Queen   ", Text: "Human,  Female, Dress, Rich, Ruler"},
	}

	for i, dp := range dataPoints {
		dataPoint := dp.(data)

		embed, err := llm.CreateEmbedding(context.Background(), []string{dataPoint.Text})
		if err != nil {
			log.Fatal(err)
		}

		dataPoint.Vec = embed[0]
		dataPoints[i] = dataPoint
	}

	// -------------------------------------------------------------------------

	for _, target := range dataPoints {
		results := vector.Similarity(target, dataPoints...)

		for _, result := range results {
			fmt.Printf("%s -> %s: %.3f%% similar\n",
				result.Target.(data).Name,
				result.DataPoint.(data).Name,
				result.Percentage)
		}
		fmt.Print("\n")
	}

	// -------------------------------------------------------------------------

	// King - Man + Woman ~= Queen

	kingSubMan := vector.Sub(dataPoints[3].Vector(), dataPoints[1].Vector())
	plusWoman := vector.Add(kingSubMan, dataPoints[2].Vector())

	result := vector.CosineSimilarity(plusWoman, dataPoints[4].Vector())
	fmt.Printf("King - Man + Woman ~= Queen similarity: %.3f%%\n", result*100)
}
