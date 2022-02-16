package main

// playing a bit with recursion, then will move to concurrency

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

// load words that have duplicat letters in them (true)?  Or skip them cause they're not great starting words (false)?
const wantWordsWithDupeLetters bool = true

// if a letter appears more than once in a word, should that letter's score be counted once (false), or for every occurence (true)?
const scoreDupeLetters bool = false

type kv struct {
	Key   string
	Value int
}

var letters = make(map[rune]int)
var words = make(map[string]int)

func initLetters() {
	for letter := 'a'; letter <= 'z'; letter++ {
		letters[letter] = 0
	}
}

func initWords(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("readLines: %s", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		if !wantWordsWithDupeLetters && hasDupeLetters(scanner.Text()) {
			// skip words that are low-value starting words..
			continue
		}

		words[scanner.Text()] = 0
	}
}

func scoreLetters() {
	for letter := range letters {
		for word := range words {
			if strings.ContainsRune(word, letter) {
				letters[letter]++
			}
		}
	}
}

func scoreWord(word string) int {
	var tested string
	var score int = 0
	for _, letter := range word {
		if string(letter) == "," {
			continue
		}
		// only score each letter once
		if !scoreDupeLetters && strings.ContainsRune(tested, letter) {
			continue
		}
		tested = tested + string(letter)
		score = score + letters[letter]
	}
	return score
}

func scoreWords() {
	for word := range words {
		words[word] = scoreWord(word)
	}
}

func sortScoredThings(scoredThings map[string]int) []kv {
	var sortedScoredThings []kv

	for k, v := range scoredThings {
		sortedScoredThings = append(sortedScoredThings, kv{k, v})
	}

	sort.Slice(sortedScoredThings, func(i, j int) bool {
		return sortedScoredThings[i].Value > sortedScoredThings[j].Value
	})

	return sortedScoredThings
}

func hasDupeLetters(word string) bool {
	for _, letter := range word {
		if letter == ',' {
			continue
		}
		if strings.Count(word, string(letter)) > 1 {
			// fmt.Printf("dupe in word %s, %c\n", word, letter)
			return true
		}
	}

	return false
}

func pruneThing(scoredThing map[string]int) map[string]int {
	sortedScoredThing := sortScoredThings(scoredThing)

	// prune to..
	topX := 10000

	// need to be within slice bounds
	currentLength := len(sortedScoredThing)
	if topX > currentLength {
		topX = currentLength
	}

	var prunedSortedScoredThing = make(map[string]int)
	for _, kv := range sortedScoredThing[0:topX] {
		prunedSortedScoredThing[kv.Key] = kv.Value
	}

	return prunedSortedScoredThing
}

func removeWord(dirtyWord string, words map[string]int) map[string]int {
	// strips a word out of the slice..
	var cleanWords = make(map[string]int)

	for word := range words {
		if word != dirtyWord {
			cleanWords[word] = words[word]
		}
	}

	return cleanWords
}

func main() {
	initLetters()
	initWords("c:\\temp\\wordlewords.txt")
	scoreLetters()

	var topX int

	// from https://github.com/tabatkins/wordle-list

	// what'd we read?
	fmt.Printf("Read %d words", len(words))

	// score and print letters
	fmt.Println("--- scored letters ---")
	for letter, score := range letters {
		fmt.Println(string(letter), score)
	}
	fmt.Println("------")

	/////////////// score words
	scoreWords()
	fmt.Printf("Scored %d words", len(words))

	////////////// top X scored words
	topX = 10
	fmt.Printf("--- top %d scored words ---\n", topX)
	// sort the scored words
	sortedScoredWords := sortScoredThings(words)

	for _, kv := range sortedScoredWords[0:topX] {
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}
	fmt.Println("------")

}
