package main

import (
	"fmt"

	"github.com/ardanlabs/vector/foundation/vector"
)

/*
	https://www.youtube.com/watch?v=72XgD322wZ8
	https://www.youtube.com/watch?v=Fuw0wv3X-0o&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=40

	"embeddings" emphasizes the notion of representing data in a meaningful and
	structured way (features).

	"vectors" refers to the numerical representation itself.
*/

type embedding struct {
	Name      string
	Authority float64
	HasTail   float64
	Rich      float64
	Gender    float64
}

func (emb embedding) Vector() []float64 {
	return []float64{
		emb.Authority,
		emb.HasTail,
		emb.Rich,
		emb.Gender,
	}
}

func (emb embedding) String() string {
	return fmt.Sprintf("%f", emb.Vector())
}

// =============================================================================

func main() {
	vectors := []vector.Embedding{
		embedding{Name: "Horse   ", Authority: 0.01, HasTail: 1.0, Rich: 0.1, Gender: +1.0},
		embedding{Name: "Man     ", Authority: 0.20, HasTail: 0.0, Rich: 0.3, Gender: -1.0},
		embedding{Name: "Woman   ", Authority: 0.20, HasTail: 0.0, Rich: 0.3, Gender: +1.0},
		embedding{Name: "King    ", Authority: 1.00, HasTail: 0.0, Rich: 1.0, Gender: -1.0},
		embedding{Name: "Queen   ", Authority: 1.00, HasTail: 0.0, Rich: 1.0, Gender: +1.0},
	}

	// -------------------------------------------------------------------------

	fmt.Print("\n")
	for _, v := range vectors {
		fmt.Printf("%s: %v\n", v.(embedding).Name, v)
	}
	fmt.Print("\n")

	for _, target := range vectors {
		results := vector.Similarity(target, vectors...)

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

	kingSubMan := vector.Sub(vectors[3].Vector(), vectors[1].Vector())
	plusWoman := vector.Add(kingSubMan, vectors[2].Vector())

	result := vector.CosineSimilarity(plusWoman, vectors[4].Vector())
	fmt.Printf("King - Man + Woman ~= Queen similarity: %.3f%%\n", result*100)
}
