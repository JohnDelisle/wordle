package main

// playing a bit with recursion, then will move to concurrency

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
)

const wantWordsWithDupeLetters bool = true

type kv struct {
	Key   string
	Value int
}

func initLetters() map[rune]int {
	var letters = make(map[rune]int)
	for letter := 'a'; letter <= 'z'; letter++ {
		letters[letter] = 0
	}
	return letters
}

func initWords(path string) (map[string]int, error) {
	var words = make(map[string]int)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
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
	return words, scanner.Err()
}

func scoreLetters(letters map[rune]int, words map[string]int) map[rune]int {
	for letter := range letters {
		for word := range words {
			if strings.ContainsRune(word, letter) {
				letters[letter]++
			}
		}
	}
	return letters
}

func scoreWord(word string, letters map[rune]int) int {
	var tested string
	var score int = 0
	for _, letter := range word {
		if string(letter) == "," {
			continue
		}
		// only score each letter once
		if strings.ContainsRune(tested, letter) {
			continue
		}
		tested = tested + string(letter)
		score = score + letters[letter]
	}
	return score
}

func scoreWords(letters map[rune]int, words map[string]int) map[string]int {
	var scoredWords = make(map[string]int)
	for word := range words {
		scoredWords[word] = scoreWord(word, letters)
	}
	return scoredWords
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

type safePairs struct {
	mu sync.Mutex
	v  map[string]int
}

func wordTwoWorker(w *sync.WaitGroup, pairs *safePairs, words map[string]int, wordOne string, wordTwos map[string]int) {
	for wordTwo := range wordTwos {
		pairWord := wordOne + "," + wordTwo
		if hasDupeLetters(pairWord) {
			continue
		}
		pairScore := words[wordOne] + words[wordTwo]
		pairs.mu.Lock()
		pairs.v[pairWord] = pairScore
		pairs.mu.Unlock()
	}

	// too much memory :(
	pairs.mu.Lock()
	pairs.v = pruneThing(pairs.v)
	pairs.mu.Unlock()
	w.Done()
}

func scorePairs(words map[string]int) map[string]int {
	pairs := safePairs{v: make(map[string]int)}
	var w sync.WaitGroup

	wordTwos := words
	for wordOne := range words {
		// remove wordOne, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"

		wordTwos = removeWord(wordOne, wordTwos)
		w.Add(1)
		go wordTwoWorker(&w, &pairs, words, wordOne, wordTwos)
	}
	return pairs.v
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
	var topX, currentLength int

	// from https://github.com/tabatkins/wordle-list
	letters := initLetters()
	words, err := initWords("c:\\temp\\wordlewords.txt")

	if err != nil {
		log.Fatalf("readLines: %s", err)
	}

	// what'd we read?
	fmt.Printf("Read %d words", len(words))

	// score and print letters
	fmt.Println("--- scored letters ---")
	scoredLetters := scoreLetters(letters, words)
	for letter, score := range scoredLetters {
		fmt.Println(string(letter), score)
	}
	fmt.Println("------")

	/////////////// score words
	scoredWords := scoreWords(letters, words)
	fmt.Printf("Scored %d words", len(scoredWords))

	////////////// top X scored words
	topX = 10
	fmt.Printf("--- top %d scored words ---\n", topX)
	// sort the scored words
	sortedScoredWords := sortScoredThings(scoredWords)

	for _, kv := range sortedScoredWords[0:topX] {
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}
	fmt.Println("------")

	/////////////// top X scored pairs
	topX = 1000
	fmt.Printf("--- top %d scored pairs ---\n", topX)
	// score the pairs
	scoredPairs := scorePairs(scoredWords)
	// sort the scored pairs
	sortedScoredPairs := sortScoredThings(scoredPairs)

	// need to be within slice bounds
	currentLength = len(sortedScoredPairs)
	if topX > currentLength {
		topX = currentLength
	}

	for _, kv := range sortedScoredPairs[0:topX] {
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}
	fmt.Println("------")

}
