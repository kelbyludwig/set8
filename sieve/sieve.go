package sieve

// https://golang.org/doc/play/sieve.go

import "math/big"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
// Generates a slice of primes until a prime exceeds n
func GeneratePrimes(n int) []int {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	result := make([]int, 0)
	prime := 1
	for prime < n {
		prime = <-ch
		result = append(result, prime)
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
	return result
}

func UniqueSmallFactors(x *big.Int) []*big.Int {
	limit := 65535
	primes := GeneratePrimes(limit)
	zero := big.NewInt(0)
	factors := make([]*big.Int, 0)
	for _, p := range primes {
		bigP := big.NewInt(int64(p))
		rem := new(big.Int).Mod(x, bigP)
		if rem.Cmp(zero) == 0 {
			factors = append(factors, bigP)
		}
	}
	return factors
}
