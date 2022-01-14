package main

import (
	"fmt"
	"math/rand"
)

func main() {
	// declare and initialize variables for operand sizes
	var r1, c1, c2 int
	r1 = 1000
	c1 = 1000
	c2 = 1000

	// declare 3 n x n matrices to mutliply and for result
	var first [1000][1000]int
	var second [1000][1000]int
	var result [1000][1000]int

	// generate random values for both matrices from 1 to 10
	for i := 0; i < r1; i++ {
		for j := 0; j < c1; j++ {
			first[i][j] = rand.Intn(10)
			second[i][j] = rand.Intn(10)
		}
	}

	// multiply both matrices
	for i := 0; i < r1; i++ {
		for j := 0; j < c2; j++ {
			for k := 0; k < c1; k++ {
				result[i][j] += first[i][k] * second[k][j]
			}
		}
	}

	// print out 10x10 submatrix of the operands and result
	fmt.Println("Matrix 1:")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print(first[i][j], " ")
		}
		fmt.Println()
	}

	fmt.Println("Matrix 2:")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print(second[i][j], " ")
		}
		fmt.Println()
	}

	fmt.Println("Result:")
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			fmt.Print(result[i][j], " ")
		}
		fmt.Println()
	}
}
