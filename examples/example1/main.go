package main

import (
	"fmt"

	"github.com/ardanlabs/vector/foundation/vector"
)

/*
	https://www.youtube.com/watch?v=72XgD322wZ8
	https://www.youtube.com/watch?v=Fuw0wv3X-0o&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=40
	https://www.youtube.com/watch?v=hQwFeIupNP0&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=41
	https://machinelearningmastery.com/what-are-word-embeddings/
	https://machinelearningmastery.com/use-word-embedding-layers-deep-learning-keras/

	"embeddings" emphasizes the notion of representing data in a meaningful and
	structured way (via features).

	"vectors" refers to the numerical representation of those features.

	Embeddings are not hand crafted like in this example. Here, we hand coded
	the features and the feature vectors. To do this at scale, this needs to
	bed automated. This is done during neural network training and you won't
	know what they features are. But it all works.
*/

type embedding struct {
	Name      string
	Authority float32 // These fields are called features.
	Animal    float32
	Human     float32
	Rich      float32
	Gender    float32
}

// Vector can convert the specified embedding into a vector.
func (emb embedding) Vector() []float32 {
	return []float32{
		emb.Authority,
		emb.Animal,
		emb.Human,
		emb.Rich,
		emb.Gender,
	}
}

// String pretty prints an embedding to a vector representation.
func (emb embedding) String() string {
	return fmt.Sprintf("%f", emb.Vector())
}

// =============================================================================

func main() {

	// Apply the feature vectors to the hand crafted embeddings.
	vectors := []vector.Embedding{
		embedding{Name: "Horse   ", Authority: 0.0, Animal: 1.0, Human: 0.0, Rich: 0.0, Gender: +1.0},
		embedding{Name: "Man     ", Authority: 0.0, Animal: 0.0, Human: 1.0, Rich: 0.0, Gender: -1.0},
		embedding{Name: "Woman   ", Authority: 0.0, Animal: 0.0, Human: 1.0, Rich: 0.0, Gender: +1.0},
		embedding{Name: "King    ", Authority: 1.0, Animal: 0.0, Human: 1.0, Rich: 1.0, Gender: -1.0},
		embedding{Name: "Queen   ", Authority: 1.0, Animal: 0.0, Human: 1.0, Rich: 1.0, Gender: +1.0},
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
