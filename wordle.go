package wordle

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

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

func scoreWords(letters map[rune]int, words map[string]int) map[string]int {
	var scoredWords = make(map[string]int)

	for word := range words {
		if hasDupeLetters(word) {
			// skip words that are low-value starting words..
			continue
		}

		var tested string
		for _, letter := range word {
			// only score each letter once
			if strings.ContainsRune(tested, letter) {
				continue
			}
			tested = tested + string(letter)

			scoredWords[word] = scoredWords[word] + letters[letter]
		}
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

func scorePairs(words map[string]int) map[string]int {
	var pairs = make(map[string]int)
	wordTwos := words

	for wordOne := range words {
		// remove wordOne, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"
		wordTwos = removeWord(wordOne, wordTwos)

		for wordTwo := range wordTwos {
			pairWord := wordOne + "," + wordTwo
			if hasDupeLetters(pairWord) {
				continue
			}
			pairScore := words[wordOne] + words[wordTwo]
			pairs[pairWord] = pairScore
		}

		// too much memory :(
		pairs = pruneThing(pairs)
	}

	return pairs
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

func scoreTrips(words map[string]int) map[string]int {
	var trips = make(map[string]int)

	wordTwos := words
	wordThrees := words

	c := 0

	for wordOne := range words {
		if hasDupeLetters(wordOne) {
			continue
		}

		// remove wordOne, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"
		wordTwos = removeWord(wordOne, wordTwos)
		fmt.Printf("wordTwos %d\n", len(wordTwos))

		for wordTwo := range wordTwos {
			pairWord := wordOne + "," + wordTwo

			if hasDupeLetters(pairWord) {
				continue
			}

			pairScore := words[wordOne] + words[wordTwo]

			// remove wordTwo, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"
			wordThrees = removeWord(wordTwo, wordThrees)
			// fmt.Printf("wordThrees %d\n", len(wordThrees))

			for wordThree := range wordThrees {

				tripsWord := pairWord + "," + wordThree

				if hasDupeLetters(tripsWord) {
					continue
				}

				tripsScore := pairScore + words[wordThree]
				trips[tripsWord] = tripsScore
				// fmt.Printf("Trips %s %d\n", tripsWord, tripsScore)
			}

		}
		c++
		fmt.Printf("Onto word %s, %d/%d\n", wordOne, c, len(words))

		// prune the trips.. too much memory used :(
		fmt.Printf("trips length: %d\n", len(trips))
		trips = pruneThing(trips)
	}

	return trips
}

func scoreQuads(words map[string]int) map[string]int {
	var quads = make(map[string]int)

	wordTwos := words
	wordThrees := words
	wordFours := words

	c := 0

	for wordOne := range words {
		if hasDupeLetters(wordOne) {
			continue
		}

		// remove wordOne, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"
		wordTwos = removeWord(wordOne, wordTwos)
		fmt.Printf("wordTwos %d\n", len(wordTwos))

		for wordTwo := range wordTwos {
			pairWord := wordOne + "," + wordTwo

			if hasDupeLetters(pairWord) {
				continue
			}

			pairScore := words[wordOne] + words[wordTwo]

			// remove wordTwo, avoid permutatioins like "cat,dog,bat" and "bat,cat,dog"
			wordThrees = removeWord(wordTwo, wordThrees)
			// fmt.Printf("wordThrees %d\n", len(wordThrees))

			for wordThree := range wordThrees {

				tripsWord := pairWord + "," + wordThree

				if hasDupeLetters(tripsWord) {
					continue
				}

				tripsScore := pairScore + words[wordThree]

				wordFours = removeWord(wordThree, wordFours)

				for wordFour := range wordFours {

					quadsWord := tripsWord + "," + wordFour

					if hasDupeLetters(quadsWord) {
						continue
					}

					quadsScore := tripsScore + words[wordFour]
					quads[quadsWord] = quadsScore
				}
			}
		}
		c++
		fmt.Printf("Onto word %s, %d/%d\n", wordOne, c, len(words))

		// prune the quads.. too much memory used :(
		fmt.Printf("trips length: %d\n", len(quads))
		quads = pruneThing(quads)
	}

	return quads
}

func main() {
	var topX int
	var currentLength int

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
	/*
		/////////////// top X scored trips
		topX = 1000

		fmt.Printf("--- top %d scored trips ---\n", topX)
		// score the trips
		scoredTrips := scoreTrips(scoredWords)
		// sort the scored pairs
		sortedScoredTrips := sortScoredThings(scoredTrips)

		// need to be within slice bounds
		currentLength = len(sortedScoredTrips)
		if topX > currentLength {
			topX = currentLength
		}

		for _, kv := range sortedScoredTrips[0:topX] {
			fmt.Printf("%s %d\n", kv.Key, kv.Value)
		}
		fmt.Println("------")

		/////////////// top X scored quads
		topX = 1000

		fmt.Printf("--- top %d scored quads ---\n", topX)
		// score the trips
		scoredQuads := scoreQuads(scoredWords)
		// sort the scored pairs
		sortedScoredQuads := sortScoredThings(scoredQuads)

		// need to be within slice bounds
		currentLength = len(sortedScoredQuads)
		if topX > currentLength {
			topX = currentLength
		}

		for _, kv := range sortedScoredQuads[0:topX] {
			fmt.Printf("%s %d\n", kv.Key, kv.Value)
		}
		fmt.Println("------") */

}
