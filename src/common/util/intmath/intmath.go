// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package intmath

func Min(i int, v ...int) int {
	min := i
	for _, j := range v {
		if j < min {
			min = j
		}
	}
	return min
}

func Max(i int, v ...int) int {
	max := i
	for _, j := range v {
		if j > max {
			max = j
		}
	}
	return max
}

func Ceil(i, j int) int {
	return (i + j - 1) / j
}
