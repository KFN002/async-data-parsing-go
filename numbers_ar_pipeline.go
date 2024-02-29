package main

func NumbersGen1(nums ...int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for _, num := range nums {
			out <- num
		}
	}()
	return out
}

func Filter1(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for num := range in {
			if num > 0 {
				out <- num
			}
		}
	}()
	return out
}

func Multiply1(in <-chan int) int {
	result := 1
	for num := range in {
		result *= num
	}
	return result
}

func MultiplyPipeline(inputNums ...[]int) int {
	numCh := make(chan int)
	go func() {
		defer close(numCh)
		for _, nums := range inputNums {
			for num := range NumbersGen1(nums...) {
				numCh <- num
			}
		}
	}()

	filteredCh := Filter1(numCh)

	result := Multiply1(filteredCh)

	return result
}
