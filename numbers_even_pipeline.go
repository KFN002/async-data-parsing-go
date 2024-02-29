package main

import (
	"bufio"
	"fmt"
	"os"
)

func NumbersGen(filename string) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		file, err := os.Open(filename)
		if err != nil {
			fmt.Println("Error opening file:", err)
			return
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			num := 0
			_, err := fmt.Sscanf(scanner.Text(), "%d", &num)
			if err == nil {
				out <- num
			}
		}
		if err := scanner.Err(); err != nil {
			fmt.Println("Error reading file:", err)
		}
	}()
	return out
}

func Filter(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			if num%2 == 0 {
				out <- num
			}
		}
	}()
	return out
}

func Sum(in <-chan int) int {
	sum := 0
	for num := range in {
		sum += num
	}
	return sum
}

func SumValuesPipeline(filename string) int {
	numbers := NumbersGen(filename)
	evens := Filter(numbers)
	result := Sum(evens)
	return result
}
