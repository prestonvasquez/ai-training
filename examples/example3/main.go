package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/ardanlabs/vector/foundation/stopwords"
	"github.com/ardanlabs/vector/foundation/vector"
	"github.com/ardanlabs/vector/foundation/word2vec"
)

/*
	https://www.youtube.com/watch?v=Q2NtCcqmIww&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=42
	http://snap.stanford.edu/data/amazon/productGraph/categoryFiles/reviews_Cell_Phones_and_Accessories_5.json.gz

	// Build The word2vec cli tooling. All the instructions are there.
	// Make the binary accessible. I put it under $GOPATH/bin
	https://github.com/maxoodf/word2vec
*/

var fileName = "reviews_Cell_Phones_and_Accessories_5"

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

	input, err := os.Open("zarf/data/" + fileName + ".json")
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	output, err := os.Create("zarf/data/" + fileName + ".txt")
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

	arg := []string{"-v", "-w", "5", "-m", "2", "-t", "12", "-f", "zarf/data/" + fileName + ".txt", "-o", "zarf/data/" + fileName + ".w2v"}

	cmd := exec.Command("w2v_trainer", arg...)
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("exec: %w", err)
	}

	fmt.Print("\n")

	return nil
}

func testModel() error {
	fmt.Println("Testing Model ...")
	fmt.Print("\n")

	w2v, err := word2vec.Load("zarf/data/"+fileName+".w2v", 300)
	if err != nil {
		return err
	}

	seq := make([]word2vec.Nearest, 10)
	w2v.Lookup("bad", seq)

	fmt.Println("Top 10 words similar to \"bad\"")
	fmt.Println(seq)
	fmt.Print("\n")

	// -------------------------------------------------------------------------

	var cheap [300]float32
	if err := w2v.VectorOf("cheap", cheap[:]); err != nil {
		return err
	}

	var inexpensive [300]float32
	if err := w2v.VectorOf("inexpensive", inexpensive[:]); err != nil {
		return err
	}

	v := vector.CosineSimilarity(cheap[:], inexpensive[:])

	fmt.Println("The similarity between the word \"cheap\" and \"inexpensive\"")
	fmt.Printf("%.3f%%\n", v*100)

	return nil
}
