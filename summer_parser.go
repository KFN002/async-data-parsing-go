package main

import (
	"fmt"
	"sync"
)

type Summer interface {
	ProcessSum(summer func(arr []int, result chan<- int), nums []int, chunkSize int) (int, error)
}

func ProcessSum(summer func(arr []int, result chan<- int), nums []int, chunkSize int) (int, error) {
	if chunkSize <= 0 {
		return 0, fmt.Errorf("chunkSize should be greater than 0")
	}
	total := 0
	var wg sync.WaitGroup
	resultCh := make(chan int, len(nums)/chunkSize+1)

	for i := 0; i < len(nums); i += chunkSize {
		end := i + chunkSize
		if end > len(nums) {
			end = len(nums)
		}

		wg.Add(1)
		go func(slice []int) {
			defer wg.Done()
			summer(slice, resultCh)
		}(nums[i:end])
	}

	go func() {
		wg.Wait()
		close(resultCh)
	}()

	for partialSum := range resultCh {
		total += partialSum
	}

	return total, nil
}

func SumChunk(arr []int, result chan<- int) {
	sum := 0
	for _, num := range arr {
		sum += num
	}
	result <- sum
}
