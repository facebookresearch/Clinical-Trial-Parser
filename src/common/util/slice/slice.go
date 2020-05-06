// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package slice

import (
	"sort"
	"strconv"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"
)

func ToIntSet(s []string) map[int]bool {
	if len(s) == 0 {
		return map[int]bool{}
	}
	u := make(map[int]bool)
	for _, a := range s {
		i, _ := strconv.Atoi(a)
		u[i] = true
	}
	return u
}

func IntSetToStringSlice(s map[int]bool) []string {
	u := make([]int, 0, len(s))
	for a := range s {
		u = append(u, a)
	}
	sort.Ints(u)
	v := make([]string, 0, len(u))
	for _, a := range u {
		v = append(v, strconv.Itoa(a))
	}
	return v
}

// setToSlice converts the set to a slice. If the set contains ints,
// the slice is sorted numerically.
func SetToSlice(s set.Set) []string {
	u := make(map[int]bool)
	for k := range s {
		if a, err := strconv.Atoi(k); err == nil {
			u[a] = true
		} else {
			return s.Slice()
		}
	}
	return IntSetToStringSlice(u)
}

func TrimSpace(v []string) {
	for i := range v {
		v[i] = strings.TrimSpace(v[i])
	}
}

func RemoveEmpty(l []string) []string {
	for i := len(l) - 1; i >= 0; i-- {
		if len(l[i]) == 0 {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}

func Dedupe(l []string) []string {
	l = RemoveEmpty(l)
	if len(l) < 2 {
		return l
	}
	sort.Strings(l)
	for i := len(l) - 1; i > 0; i-- {
		if l[i-1] == l[i] {
			l = append(l[:i], l[i+1:]...)
		}
	}
	return l
}
