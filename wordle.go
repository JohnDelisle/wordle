package main

// playing a bit with recursion, then will move to concurrency

// NOPE, not working out as planned, nicer looking (?) than loops within loops.. but can't shrink candidate words, run-time is brutal

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

// how many top scores to track?  ie how many starting word combos to keep
const keepTopXScores int = 100

type kv struct {
	Key   string
	Value int
}

var letters = make(map[string]int)
var allWords = make(map[string]int)

func initLetters() {
	for letter := 'a'; letter <= 'z'; letter++ {
		letters[string(letter)] = 0
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

		allWords[scanner.Text()] = 0
	}
}

func scoreLetters() {
	for letter := range letters {
		for word := range allWords {
			if strings.Contains(word, letter) {
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
		score = score + letters[string(letter)]
	}
	return score
}

func scoreWords() {
	for word := range allWords {
		allWords[word] = scoreWord(word)
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

func findStarters(startingCombos map[string]int, shrinkingWords map[string]int, depth int) map[string]int {
	if depth <= 0 {
		fmt.Println("depth 0, returning startingCombos")
		for startingCombo := range startingCombos {
			startingCombos[startingCombo] = scoreWord(startingCombo)
		}
		return startingCombos
	}

	// to hold the new longer combos as we iterate recursively
	var newStartingCombos = make(map[string]int)

	// first pass, populate with "words"
	if len(startingCombos) == 0 {
		fmt.Println("first pass, populating startingCombos")
		for word := range shrinkingWords {
			newStartingCombos[word] = 0
		}
	} else {
		for startingCombo := range startingCombos {
			for word := range shrinkingWords {
				newStartingCombo := startingCombo + "," + word
				if uniqueStartingCombo(newStartingCombo, newStartingCombos) {
					newStartingCombos[newStartingCombo] = 0
				}
			}
		}
	}

	return findStarters(newStartingCombos, shrinkingWords, depth-1)
}

func uniqueStartingCombo(newStartingCombo string, newStartingCombos map[string]int) bool {
	var foundComboMatch bool
	splitWords := strings.Split(newStartingCombo, ",")
	for startingCombo := range newStartingCombos {
		foundComboMatch = true // assume match until proven otherwise
		for i := 0; i < len(splitWords); i++ {
			if !strings.Contains(startingCombo, splitWords[i]) {
				foundComboMatch = false
			}
		}
		if foundComboMatch {
			return false
		}
	}
	return true
}

func buildStartingCombos() map[string]int {
	// four layers deep
	var startingCombos = make(map[string]int)
	var wordTwos, wordThrees, wordFours map[string]int
	//var wordOnes, wordTwos, wordThrees, wordFours map[string]int

	var topX []int
	var keep bool

	wordTwos = allWords
	for wordOne := range allWords {
		wordTwos = removeWord(wordOne, wordTwos)

		wordThrees = wordTwos
		for wordTwo := range wordTwos {
			wordThrees = removeWord(wordTwo, wordThrees)

			wordFours = wordThrees
			for wordThree := range wordThrees {
				wordFours = removeWord(wordThree, wordFours)

				for wordFour := range wordFours {

					// TODO - make this loop run concurrently, performance is crap (can't pre-score the words with new scoring idea (scoring the combo, rather than summing the word scores))

					startingCombo := wordOne + "," + wordTwo + "," + wordThree + "," + wordFour
					// startingComboScore := allWords[wordOne] + allWords[wordTwo] + allWords[wordThree] + allWords[wordFour]
					startingComboScore := scoreWord(startingCombo)

					keep, topX = worthKeeping(startingComboScore, topX)
					if keep {
						startingCombos[startingCombo] = startingComboScore
					}
				}
			}
		}
	}
	return startingCombos
}

func worthKeeping(score int, topScores []int) (bool, []int) {
	// we want to populate at least keepTopXScores worth of scores..
	if len(topScores) <= keepTopXScores {
		topScores = append(topScores, score)
		return true, topScores
	}

	// sort topX scores
	sort.Ints(topScores)

	// walk sorted slice biggest to smallest value - is our score bigger?
	for i := len(topScores) - 1; i == 0; i-- {
		if score > topScores[i] {
			topScores[i] = score
			return true, topScores
		}
	}
	return false, topScores
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

func printSortedScoredThing(sortedScoredThing []kv, topX int) {
	if topX > len(sortedScoredThing) || topX == 0 {
		topX = len(sortedScoredThing)
	}

	for _, kv := range sortedScoredThing[:topX] {
		fmt.Printf("%s %d\n", kv.Key, kv.Value)
	}
}

func main() {
	initLetters()
	initWords("c:\\temp\\wordlewords.txt.short")
	scoreLetters()

	var topX int

	// from https://github.com/tabatkins/wordle-list

	// what'd we read?
	fmt.Printf("Read %d words", len(allWords))

	// score and print letters
	fmt.Println("--- scored letters ---")
	sortedScoredLettes := sortScoredThings(letters)
	printSortedScoredThing(sortedScoredLettes, 0)
	fmt.Println("------")

	/////////////// score words
	scoreWords()
	fmt.Printf("Scored %d words", len(allWords))

	////////////// top X scored words
	topX = 10
	fmt.Printf("--- top %d scored words ---\n", topX)
	// sort the scored words
	sortedScoredWords := sortScoredThings(allWords)
	printSortedScoredThing(sortedScoredWords, topX)
	fmt.Println("------")

	////////////// starting combos
	topX = 10
	fmt.Printf("--- top %d starting combos ---\n", topX)
	// build our word combos
	startingCombos := buildStartingCombos()
	// sort the scored words
	sortedStartingCombos := sortScoredThings(startingCombos)
	printSortedScoredThing(sortedStartingCombos, 0)
	fmt.Println("------")

}
