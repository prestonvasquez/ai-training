// This example show you how to use MongoDB and Ollama to create a proper vector
// embedding database of the Ultimate Go Notebook. With this vector database,
// you will be able to query for content that has a strong similarity to your
// question.
//
// The book has already been pre-processed into chunks based on the books TOC.
// For chunks over 500 words, those chunks have been chunked again into 250
// blocks. The code will create a vector embedding for each chunk.
// That data can be found under `zarf/data/book.chunks`.
//
// The original version of the book in text format has been retained. The program
// to clean that document into chunks can be found under `cmd/cleaner`. You can
// run that program using `make clean-data`. This is here if you want to play
// with your own chunking. How you chunk the data is critical to accuracy.
//
// # Running the example:
//
//   $ make example5
//
// # This requires running the following command:
//
//   $ make dev-up // This starts the mongodb and ollama service in docker compose.

package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/tmc/langchaingo/llms/ollama"
)

/*
	Builder Process
	book.chunks -> Vectorize -> Store in MongoDB With Metadata

	Chat Process
	Question -> Vectorize -> Query MongoDB -> Pass top 5 Chunks to Model with
	                                          question

	NOTE: You must run `make dev-up` to run this example.
*/

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	if err := vectorize(); err != nil {
		return fmt.Errorf("vectorize: %w", err)
	}

	return nil
}

func vectorize() error {
	llm, err := ollama.New(ollama.WithModel("mxbai-embed-large"))
	if err != nil {
		return fmt.Errorf("ollama: %w", err)
	}

	input, err := os.Open("zarf/data/book.chunks")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer input.Close()

	var counter int

	fmt.Print("\033[s")

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		counter++

		v := scanner.Text()

		fmt.Print("\033[u\033[K")
		fmt.Printf("Vectorizing Data: %d of 341", counter)

		_, err := llm.CreateEmbedding(context.Background(), []string{v})
		if err != nil {
			return fmt.Errorf("create embedding: %w", err)
		}
	}

	fmt.Print("\n")

	return nil
}
