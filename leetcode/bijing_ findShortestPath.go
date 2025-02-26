package leetcode

import (
	"fmt"
)

func findShortestPath(grid [][]int) (int, []string) {
	rows := len(grid)
	cols := len(grid[0])
	dp := make([][]int, rows)
	for i := range dp {
		dp[i] = make([]int, cols)
	}

	dp[0][0] = grid[0][0]

	// 填充第一行
	for j := 1; j < cols; j++ {
		dp[0][j] = dp[0][j-1] + grid[0][j]
	}

	// 填充第一列
	for i := 1; i < rows; i++ {
		dp[i][0] = dp[i-1][0] + grid[i][0]
	}

	// 填充其他格子
	for i := 1; i < rows; i++ {
		for j := 1; j < cols; j++ {
			dp[i][j] = min(dp[i-1][j], dp[i][j-1]) + grid[i][j]
		}
	}

	// 回溯路径
	path := []string{}
	i, j := rows-1, cols-1
	for i > 0 || j > 0 {
		path = append([]string{posToStr(i, j)}, path...)
		if i == 0 {
			j--
		} else if j == 0 {
			i--
		} else {
			if dp[i-1][j] < dp[i][j-1] {
				i--
			} else {
				j--
			}
		}
	}
	path = append([]string{"A1"}, path...)

	return dp[rows-1][cols-1], path
}

func posToStr(i, j int) string {
	row := i + 1
	col := 'A' + j
	return fmt.Sprintf("%c%d", col, row)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
