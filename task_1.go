// Task: Sum of Numbers
// Write a Go function that takes a slice of integers as input and returns the sum of all the numbers. If the slice is empty, the function should return 0.
// [Optional]: Write a test for your function


package main

import "fmt"

func Sum(nums []int) int {
    sum := 0
    for i:=0; i<len(nums); i++ {
	 	sum += nums[i]
	}
    return sum
}

func main() {
    fmt.Println(Sum([]int{1, 2, 3, 4, 5}))
    fmt.Println(Sum([]int{}))
}