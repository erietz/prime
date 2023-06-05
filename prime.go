package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"sync"
)

var (
	low         = flag.Int("lo", 0, "Lower bound of range")
	high        = flag.Int("hi", 100, "Upper bound of range")
	nPartitions = flag.Int("nPartitions", 1, "Number of goroutines to divide the work into")
)

type Partition struct {
	low  int
	high int
}

// returns true if n is prime, else false
func isPrime(n int) bool {
	for i := 2; i <= n/2; i++ {
		if n%i == 0 {
			return false
		}
	}

	return true
}

// returns all of the primes between lo and hi
func findPrimes(low, high int) <-chan int {
	primes := make(chan int)
	go func() {
		for i := low; i < high; i++ {
			if isPrime(i) {
				primes <- i
			}
		}
		close(primes)
	}()
	return primes
}

// func combine1(chan1, chan2 chan int) chan int {
// 	c := make(chan int)
// 	go func() {
// 		ok1, ok2 := true, true
// 		for {
// 			var i int
// 			select {
// 			case i, ok1 = <-chan1:
// 				if ok1 {
// 					c <- i
// 				}
// 				if !ok1 && !ok2 {
// 					close(c)
// 					return
// 				}
// 			case i, ok2 = <-chan2:
// 				if ok2 {
// 					c <- i
// 				}
// 				if !ok1 && !ok2 {
// 					close(c)
// 					return
// 				}
// 			}
// 		}
// 	}()
// 	return c
// }

// func merge(channels ...chan int) chan int {
// 	results := make(chan int)

// 	go func() {
// 		wg := sync.WaitGroup{}

// 		for _, channel := range channels {
// 			wg.Add(1)
// 			go func(ch chan int) {
// 				defer wg.Done()
// 				for result := range ch {
// 					results <- result
// 				}
// 			}(channel)
// 		}

// 		wg.Wait()
// 		close(results)
// 	}()

// 	return results
// }

// merges all of the channels into a single combined result channel.
func merge(channels ...<-chan int) <-chan int {
	results := make(chan int)

	wg := sync.WaitGroup{}
	wg.Add(len(channels))

	for _, channel := range channels {
		go func(ch <-chan int) {
			defer wg.Done()
			for i := range ch {
				results <- i
			}
		}(channel)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func checkFlags() {
	if *low > *high {
		log.Fatal("low must be less than or equal to high")
	}

	if *nPartitions > (*high - *low) {
		log.Fatal("cannot have more partitions than (high - low)")
	}
}

func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func partition(low, high, nPartitions int) []Partition {
	partitions := []Partition{}
	chunkSize := (high - low) / nPartitions
	l, h := low, low+chunkSize
	for i := 0; i < nPartitions; i++ {
		partitions = append(partitions, Partition{
			low:  l,
			high: h,
		})
		l += chunkSize
		h = min(h+chunkSize, high)
	}

	return partitions
}

func main() {
	flag.Parse()
	checkFlags()
	partitions := partition(*low, *high, *nPartitions)

	var routines []<-chan int
	for _, part := range partitions {
		fmt.Println(part.low, part.high)
		routines = append(routines, findPrimes(part.low, part.high))
	}

	results := merge(routines...)

	sortedResults := []int{}
	for result := range results {
		sortedResults = append(sortedResults, result)
	}

	sort.Ints(sortedResults)
	fmt.Println(sortedResults)
}
