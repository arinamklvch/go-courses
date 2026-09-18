package main

import "fmt"

func twoSum(nums []int, target int) []int {
	m := map[int]int{}
	res := make([]int, 0, 2)
	for i := 0; i < len(nums); i++ {
		delta := target - nums[i]
		_, ok := m[delta]
		if ok {
			res = append(res, i, m[delta])
			break
		}
		m[nums[i]] = i
	}
	return res
}

func main() {
	nums := []int{2, 7, 11, 15}
	target := 9
	fmt.Println(twoSum(nums, target))
}
