package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/ardanlabs/vector/foundation/stopwords"
)

// https://www.youtube.com/watch?v=Q2NtCcqmIww&list=PLeo1K3hjS3uu7CxAacxVndI4bE_o3BDtO&index=42
// http://snap.stanford.edu/data/amazon/productGraph/categoryFiles/reviews_Cell_Phones_and_Accessories_5.json.gz

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	data, err := readData()
	if err != nil {
		return fmt.Errorf("readData: %w", err)
	}

	fmt.Println(data[0])
	fmt.Println(data[1])

	return nil
}

func readData() ([]string, error) {
	/*
		{
		  "reviewerID": "A30TL5EWN6DFXT",
		  "asin": "120401325X",
		  "reviewerName": "christina",
		  "helpful": [
		    0,
		    0
		  ],
		  "reviewText": "They look good and stick good! I just don't like the rounded shape because I was always bumping it and Siri kept popping up and it was irritating. I just won't buy a product like this again",
		  "overall": 4,
		  "summary": "Looks Good",
		  "unixReviewTime": 1400630400,
		  "reviewTime": "05 21, 2014"
		}
	*/

	type document struct {
		ReviewText string
	}

	f, err := os.Open("zarf/data/reviews_Cell_Phones_and_Accessories_5.json")
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}

	var results []string
	var counter int

	fmt.Print("\033[s")

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		s := scanner.Text()

		var d document
		err := json.Unmarshal([]byte(s), &d)
		if err != nil {
			return nil, fmt.Errorf("unmarshal: %w", err)
		}

		v := stopwords.Remove(d.ReviewText)
		results = append(results, v)

		counter++

		fmt.Print("\033[u\033[K")
		fmt.Printf("Reading: %d", counter)
	}

	fmt.Print("\n")

	return results, nil
}
