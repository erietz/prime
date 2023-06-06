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
		routines = append(routines, findPrimes(part.low, part.high))
	}

	primes := []int{}
	for prime := range merge(routines...) {
		primes = append(primes, prime)
	}

	sort.Ints(primes)
	fmt.Println(primes)
}
