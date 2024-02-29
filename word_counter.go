package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

type WordCounter struct {
	wordsCount map[string]int
	mu         sync.Mutex
}

type CounterWorker interface {
	ProcessFiles(files ...string) error
	ProcessReader(r io.Reader) error
}

func NewWordCounter() *WordCounter {
	return &WordCounter{
		wordsCount: make(map[string]int),
	}
}

func (wc *WordCounter) ProcessFiles(files ...string) error {
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)
		go func(filename string) {
			defer wg.Done()
			err := wc.processFile(filename)
			if err != nil {
				fmt.Printf("Error processing file %s: %v\n", filename, err)
			}
		}(file)
	}

	wg.Wait()

	return nil
}

func (wc *WordCounter) ProcessReader(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		words := strings.Fields(line)

		wc.mu.Lock()
		for _, word := range words {
			lowercaseWord := strings.ToLower(word)
			wc.wordsCount[lowercaseWord]++
		}
		wc.mu.Unlock()
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (wc *WordCounter) processFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return wc.ProcessReader(file)
}
