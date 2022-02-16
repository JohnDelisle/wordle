package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
)

func initWords(path string) []string {
	var words []string

	file, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		words = append(words, scanner.Text())
	}

	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// what'd we read?
	fmt.Printf("Read %d words\n\n\n", len(words))

	return words
}

func getThing(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s", prompt)
	thing, _ := reader.ReadString('\n')
	thing = strings.Replace(thing, "\r", "", -1)
	thing = strings.Replace(thing, "\n", "", -1)

	return thing
}

func getKnown() string {
	known := getThing("What letters do we know (. for wild): ")
	return known
}

func getPositiveHits() string {
	positiveHits := getThing("What letters do we have positive hits on: ")
	return positiveHits
}

func getNegativeHits() string {
	negativeHits := getThing("What letters do we have negative hits on: ")
	return negativeHits
}

func filterKnown(known string, words []string) []string {
	fmt.Println("Filtering for known letters...")

	var filtered []string
	for _, word := range words {
		result, err := regexp.MatchString(known, word)

		if err != nil {
			log.Fatalf("Regex error: %s", err)
		}

		if result {
			filtered = append(filtered, word)
		}
	}

	printWords(filtered)
	return filtered
}

func filterPositiveHits(positiveHits string, words []string) []string {
	fmt.Printf("Filtering for positive hit letters %q...\n", positiveHits)

	var filtered []string
	for _, word := range words {
		var positiveHit bool = true

		for _, letter := range positiveHits {
			fmt.Printf("Does %q have a %q... ", word, letter)
			result, _ := regexp.MatchString(string(letter), word)
			fmt.Printf("%v\n", result)

			// word needs to contain all positive hit letters
			positiveHit = positiveHit && result
		}

		if positiveHit {
			filtered = append(filtered, word)
		}
	}

	printWords(filtered)
	return filtered
}

func filterNegativeHits(negativeHits string, words []string) []string {
	fmt.Printf("Filtering out negative hit letters %q...\n", negativeHits)

	var filtered []string
	for _, word := range words {
		var negativeHit bool

		for _, letter := range negativeHits {
			fmt.Printf("Does %q have a %q... ", word, letter)
			result, _ := regexp.MatchString(string(letter), word)
			fmt.Printf("%v\n", result)

			// word needs to contain all positive hit letters
			negativeHit = negativeHit || result
		}

		// negativeHit == false is good
		if !negativeHit {
			filtered = append(filtered, word)
		}
	}

	printWords(filtered)
	return filtered
}

func printWords(words []string) {
	fmt.Printf("%d candidate words:\n", len(words))
	sort.Strings(words)
	for _, word := range words {
		fmt.Println(word)
	}
	fmt.Println()
}

func main() {
	candidates := initWords("c:\\temp\\wordlewords.txt")

	known := getKnown()
	positiveHits := getPositiveHits()
	negativeHits := getNegativeHits()
	fmt.Printf("Using %q known letters, %q positive hits, and %q negative hits\n\n\n", known, positiveHits, negativeHits)

	candidates = filterKnown(known, candidates)
	candidates = filterPositiveHits(positiveHits, candidates)
	candidates = filterNegativeHits(negativeHits, candidates)

	printWords(candidates)
}
