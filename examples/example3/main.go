package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/ardanlabs/vector/foundation/stopwords"
	"github.com/ardanlabs/vector/foundation/vector"
	"github.com/ardanlabs/vector/foundation/word2vec"
)

/*
	https://www.youtube.com/watch?v=Q2NtCcqmIww&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=42
	http://snap.stanford.edu/data/amazon/productGraph/categoryFiles/reviews_Cell_Phones_and_Accessories_5.json.gz

	NOTE: You must run `make download-data` to get the data file need to run
	      this example.
*/

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := cleanData(); err != nil {
		return fmt.Errorf("cleanData: %w", err)
	}

	if err := trainModel(); err != nil {
		return fmt.Errorf("trainModel: %w", err)
	}

	if err := testModel(); err != nil {
		return fmt.Errorf("trainModel: %w", err)
	}

	return nil
}

func cleanData() error {
	type document struct {
		ReviewText string
	}

	input, err := os.Open("zarf/data/example3.json")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	output, err := os.Create("zarf/data/example3.words")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	var counter int

	fmt.Print("\033[s")

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		s := scanner.Text()

		var d document
		err := json.Unmarshal([]byte(s), &d)
		if err != nil {
			return fmt.Errorf("unmarshal: %w", err)
		}

		v := stopwords.Remove(d.ReviewText)

		output.WriteString(v)
		output.WriteString("\n")

		counter++

		fmt.Print("\033[u\033[K")
		fmt.Printf("Reading/Cleaning Data: %d", counter)
	}

	fmt.Print("\n")

	return nil
}

func trainModel() error {
	fmt.Println("Training Model ...")
	fmt.Print("\n")

	config := word2vec.Config{
		Corpus: word2vec.ConfigCorpus{
			InputFile: "zarf/data/example3.words",
			Tokenizer: " \n,.-!?:;/\"#$%&'()*+<=>@[]\\^_`{|}~\t\v\f\r",
			Sequencer: ".\n?!",
		},
		Vector: word2vec.ConfigWordVector{
			Vector:    300,
			Window:    5,
			Threshold: 1e-3,
			Frequency: 5,
		},
		Learning: word2vec.ConfigLearning{
			Epoch: 10,
			Rate:  0.05,
		},
		UseSkipGram:            false,
		UseCBOW:                true,
		UseNegativeSampling:    true,
		UseHierarchicalSoftMax: false,
		SizeNegativeSampling:   5,
		Threads:                runtime.GOMAXPROCS(0),
		Verbose:                true,
		Output:                 "zarf/data/example3.model",
	}

	if err := word2vec.Train(config); err != nil {
		return fmt.Errorf("train: %w", err)
	}

	fmt.Print("\n")

	return nil
}

func testModel() error {
	fmt.Println("Testing Model ...")
	fmt.Print("\n")

	w2v, err := word2vec.Load("zarf/data/example3.model", 300)
	if err != nil {
		return err
	}

	seq := make([]word2vec.Nearest, 10)
	w2v.Lookup("bad", seq)

	fmt.Println("Top 10 words similar to \"bad\"")
	fmt.Println(seq)
	fmt.Print("\n")

	// -------------------------------------------------------------------------

	words := []string{"terrible", "horrible", "price", "battery", "great", "nice"}

	for i := 0; i < len(words); i = i + 2 {
		var word1 [300]float32
		if err := w2v.VectorOf(words[i], word1[:]); err != nil {
			return err
		}

		var word2 [300]float32
		if err := w2v.VectorOf(words[i+1], word2[:]); err != nil {
			return err
		}

		v := vector.CosineSimilarity(word1[:], word2[:])

		fmt.Printf("The cosine similarity between the word %q and %q: %.3f%%\n", words[i], words[i+1], v*100)
	}

	return nil
}
