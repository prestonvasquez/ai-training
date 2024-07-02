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
	TEXT -> 400 Word Chunks -> Vectorize -> Store in MongoDB With Metadata

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

	input, err := os.Open("zarf/data/book.clean")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer input.Close()

	var counter int
	var idx int
	buffer := make([]string, 10)

	fmt.Print("\033[s")

	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		counter++

		buffer[idx] = scanner.Text()
		idx++

		if idx == 10 {
			fmt.Print("\033[u\033[K")
			fmt.Printf("Vectorizing Data: %d of 6891", counter)

			_, err := llm.CreateEmbedding(context.Background(), buffer)
			if err != nil {
				return fmt.Errorf("create embedding: %w", err)
			}

			for i := range buffer {
				buffer[i] = ""
			}
			idx = 0
		}
	}

	if idx > 0 {
		fmt.Print("\033[u\033[K")
		fmt.Printf("Vectorizing Data: %d of 6891", counter)

		_, err := llm.CreateEmbedding(context.Background(), buffer)
		if err != nil {
			return fmt.Errorf("create embedding: %w", err)
		}
	}

	fmt.Print("\n")

	return nil
}
