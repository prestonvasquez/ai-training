// This example shows you what a vector and embedding is by hand crafting
// a relationship of data. It also shows you how cosine similarity works between
// different vectors.
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

	Fields are not hand crafted like in this example. Here, we hand coded
	the features and the values. To do this at scale, this needs to be
	automated. This is done during neural network training and you won't
	know what the features are. But it all works.
*/

type data struct {
	Name      string
	Authority float32 // These fields are called features.
	Animal    float32
	Human     float32
	Rich      float32
	Gender    float32
}

// Vector can convert the specified data into a vector.
func (d data) Vector() []float32 {
	return []float32{
		d.Authority,
		d.Animal,
		d.Human,
		d.Rich,
		d.Gender,
	}
}

// String pretty prints an embedding to a vector representation.
func (d data) String() string {
	return fmt.Sprintf("%f", d.Vector())
}

// =============================================================================

func main() {

	// Apply the feature dataPoints to the hand crafted embeddings.
	dataPoints := []vector.Data{
		data{Name: "Horse   ", Authority: 0.0, Animal: 1.0, Human: 0.0, Rich: 0.0, Gender: +1.0},
		data{Name: "Man     ", Authority: 0.0, Animal: 0.0, Human: 1.0, Rich: 0.0, Gender: -1.0},
		data{Name: "Woman   ", Authority: 0.0, Animal: 0.0, Human: 1.0, Rich: 0.0, Gender: +1.0},
		data{Name: "King    ", Authority: 1.0, Animal: 0.0, Human: 1.0, Rich: 1.0, Gender: -1.0},
		data{Name: "Queen   ", Authority: 1.0, Animal: 0.0, Human: 1.0, Rich: 1.0, Gender: +1.0},
	}

	// -------------------------------------------------------------------------

	fmt.Print("\n")
	for _, v := range dataPoints {
		fmt.Printf("%s: %v\n", v.(data).Name, v)
	}
	fmt.Print("\n")

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
