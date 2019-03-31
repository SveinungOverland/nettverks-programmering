package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"sort"
	"time"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	const threads = 5
	const nrNumbers = 10000000

	slices := createNumberSlices(threads, nrNumbers)

	fmt.Println("Done making slices, assigning to threads")

	start := time.Now()
	results := setupPrime(slices)

	fmt.Print("Waiting for result nr: ")
	mergedResult := []int{2}
	for i := 0; i < threads; i++ {
		fmt.Printf("%d, ", i)
		res := <-results
		mergedResult = append(mergedResult, res...)
	}

	timeToResult := time.Now().Sub(start)

	sort.Ints(mergedResult)
	fmt.Println()
	if len(mergedResult) < 200 {
		fmt.Println(mergedResult)
	} else {
		fmt.Println(mergedResult[:200])
	}

	timeToSort := time.Now().Sub(start)

	fmt.Printf("Time spent to result: %fms \n", float64(timeToResult.Nanoseconds())*1e-6)
	fmt.Printf("Time spent to sort: %fms \n", float64(timeToSort.Nanoseconds())*1e-6)

	fmt.Println("Checking calculated primes with naive division test")
	fmt.Printf("How many prime numbers %d\n", len(mergedResult))
	for _, num := range mergedResult[:10000] {
		fmt.Printf("\r%d", num)
		for i := 2; i < num; i++ {
			if num%i == 0 {
				fmt.Printf("\nNumber %d is not prime\n", num)
			}
		}
	}
	fmt.Println("\ndone")

}

func createNumberSlices(threads, nrNumbers int) [][]int {
	estSize := int(((nrNumbers / 2) + 1) / threads)
	output := make([][]int, threads)
	outputHead := 0
	outputHeightHead := 0
	for i := range output {
		output[i] = make([]int, estSize)
	}
	for i := 3; i < nrNumbers; i += 2 {
		if outputHead == threads {
			outputHead = 0
			outputHeightHead++
		}
		output[outputHead][outputHeightHead] = i
		outputHead++
	}
	return output
}

func setupPrime(slices [][]int) chan []int {
	results := make(chan []int, len(slices))
	for i := 0; i < len(slices); i++ {
		go worker(slices[i], results)
	}
	return results
}

func worker(numbers []int, results chan<- []int) {
	/*
		Super shady code....
		sets filteredOutput to be an empty slice at the same position as numbers
		... uses the fact that a slice shares the same backing array and capacity
		as the original, so the storage is reused for the filtered slice. The
		original slice is rip....
	*/
	filteredOutput := numbers[:0]
	for _, x := range numbers {
		//fmt.Printf("Checking number: %d \n", x)
		if checkPrime(x) {
			filteredOutput = append(filteredOutput, x)
		}
	}
	results <- filteredOutput
}

func checkPrime(number int) bool {
	if number == 3 || number == 5 || number == 7 {
		return true
	}
	if number%3 == 0 || number%5 == 0 || number%7 == 0 {
		return false
	}
	for i := 11; i*i < number+1; i += 2 {
		if number%i == 0 {
			return false
		}
	}
	return true
}
