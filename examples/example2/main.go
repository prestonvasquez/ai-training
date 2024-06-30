package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ardanlabs/vector/foundation/vector"
	"github.com/tmc/langchaingo/llms/ollama"
)

type embedding struct {
	Name string
	Text string
	Vec  []float32
}

// Vector can convert the specified embedding into a vector.
func (emb embedding) Vector() []float32 {
	return emb.Vec
}

// =============================================================================

func main() {
	llm, err := ollama.New(ollama.WithModel("mxbai-embed-large"))
	if err != nil {
		log.Fatal(err)
	}

	// -------------------------------------------------------------------------

	// Apply the feature vectors to the hand crafted embeddings.
	embeddings := []vector.Embedding{
		embedding{Name: "Horse   ", Text: "Animal, Female"},
		embedding{Name: "Man     ", Text: "Human,  Male,   Pants, Poor, Worker"},
		embedding{Name: "Woman   ", Text: "Human,  Female, Dress, Poor, Worker"},
		embedding{Name: "King    ", Text: "Human,  Male,   Pants, Rich, Ruler"},
		embedding{Name: "Queen   ", Text: "Human,  Female, Dress, Rich, Ruler"},
	}

	for i, emb := range embeddings {
		e := emb.(embedding)

		embed, err := llm.CreateEmbedding(context.Background(), []string{e.Text})
		if err != nil {
			log.Fatal(err)
		}

		e.Vec = embed[0]

		embeddings[i] = e
	}

	// -------------------------------------------------------------------------

	for _, target := range embeddings {
		results := vector.Similarity(target, embeddings...)

		for _, result := range results {
			fmt.Printf("%s -> %s: %.3f%% similar\n",
				result.Target.(embedding).Name,
				result.Record.(embedding).Name,
				result.Percentage)
		}
		fmt.Print("\n")
	}

	// -------------------------------------------------------------------------

	// King - Man + Woman ~= Queen

	kingSubMan := vector.Sub(embeddings[3].Vector(), embeddings[1].Vector())
	plusWoman := vector.Add(kingSubMan, embeddings[2].Vector())

	result := vector.CosineSimilarity(plusWoman, embeddings[4].Vector())
	fmt.Printf("King - Man + Woman ~= Queen similarity: %.3f%%\n", result*100)
}
