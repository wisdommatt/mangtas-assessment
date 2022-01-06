package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type wordCount struct {
	Word  string `json:"word"`
	Count int    `json:"count"`
}

var port = "8080"

func main() {
	serveMux := http.NewServeMux()
	serveMux.HandleFunc("/", mostUsedWordsHandler)

	log.Printf("app running on port: %s", port)
	log.Fatal(http.ListenAndServe(":"+port, serveMux))
}

func mostUsedWordsHandler(rw http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(rw, "only POST request accepted", http.StatusBadRequest)
		return
	}
	var text string
	err := json.NewDecoder(r.Body).Decode(&text)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	wordCounts := extractWordsCount(text)
	if len(wordCounts) < 10 {
		http.Error(rw, "text does not contain upto 10 unique words", http.StatusBadRequest)
		return
	}
	json.NewEncoder(rw).Encode(wordCounts[0:10])
}

// extractWordsCount returns occurence count for all the words
// in the text in ascending order.
func extractWordsCount(text string) []wordCount {
	wordsMap := map[string]int{}
	textArray := strings.Split(text, " ")
	for _, word := range textArray {
		word = strings.TrimSpace(word)
		word = strings.ToLower(word)             // converting the word to lowercase to avoid duplicate entries for diffent cases
		word = strings.ReplaceAll(word, ".", "") // removing dots from the word
		if _, ok := wordsMap[word]; !ok {
			wordsMap[word] = 1
			continue
		}
		wordsMap[word]++
	}
	wordCounts := []wordCount{}
	for word, count := range wordsMap {
		wordCounts = append(wordCounts, wordCount{
			Word:  word,
			Count: count,
		})
		i := len(wordCounts) - 1
		if i < 1 {
			continue
		}
		for i > 0 && wordCounts[i-1].Count < wordCounts[i].Count {
			previousWordCount := wordCounts[i-1]
			wordCounts[i-1] = wordCount{
				Word:  word,
				Count: count,
			}
			wordCounts[i] = previousWordCount
			i--
		}
	}
	return wordCounts
}
