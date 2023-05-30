package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
)

// begin splits the work between goroutines and initialises
func begin(done chan interface{}, filename string, length int) {
	wordFile, err := ioutil.ReadFile(filename)

	if err != nil {
		log.Fatal(err)
	}

	numFinders := runtime.NumCPU() // number of goroutines we can create

	split := len(wordFile) / numFinders // how many bytes we will give to each goroutine

	newFileName := fmt.Sprintf("length_%d", length)
	// create the file to write to
	writeTo, err := os.Create(newFileName)
	defer writeTo.Close()
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup

	// fan out
	for startPoint := 0; startPoint < len(wordFile); startPoint += split {
		for j := startPoint; ; j++ {
			if wordFile[startPoint] == 10 { // if y equals a space then continue with y loop
				break
			}
			startPoint++ // Don't want to cut the slice half way through a word so find a space
		}
		wg.Add(1)
		parseWords(&wg, done, wordFile, startPoint, length, writeTo) // send a goroutine to write the given number of bytes to a file in string form
	}

	wg.Wait() // Wait until all goroutines have completed
}

// parseWords concurrently parses the file given and returns a slice of words of the given length
func parseWords(wg *sync.WaitGroup, done <-chan interface{}, values []byte, startPoint int, length int, file *os.File) {
	defer wg.Done()

	wordStream := make(chan []string)
	go func() {
		defer close(wordStream)
		// check if the end point is greater than the length of values and decrease the endpoint
		endPoint := startPoint + len(values)/runtime.NumCPU()
		for j := endPoint; j > startPoint; j-- {
			if endPoint > len(values) {
				endPoint--
			} else if endPoint == len(values) {
				break
			}
		}

		var wordSlice []string

		// add a word at each whitespace
		spaceCount := 0
		for i := startPoint; i < endPoint; i++ {
			if values[i] == 10 { // if byte is whitespace then append the word
				word := strings.TrimSpace(string(values[i-spaceCount : i]))
				if len(word) == length {
					wordSlice = append(wordSlice, word)
				}
				spaceCount = 0
			}
			spaceCount++
		}
		select {
		case <-done:
			return
		case wordStream <- wordSlice:
		}
	}()
	for _, v := range <-wordStream {
		v = v + " "
		fmt.Fprintf(file, v)
	}

}
