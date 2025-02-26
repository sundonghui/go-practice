package leetcode

import (
	"fmt"
	"testing"
)

func Test_numberNPrime(t *testing.T) {
	numberNPrime(22222223, 1000)
}

func Test_findShortestPath(t *testing.T) {
	grid := [][]int{
		{2, 4, 6, 8, 7, 9},
		{3, 8, 9, 1, 4, 5},
		{6, 9, 12, 2, 15, 11},
		{8, 7, 18, 7, 13, 2},
		{7, 2, 10, 4, 3, 7},
		{9, 5, 6, 17, 5, 1},
	}

	totalTime, path := findShortestPath(grid)
	fmt.Println("最短用时:", totalTime)
	fmt.Println("路径:", path)
}
