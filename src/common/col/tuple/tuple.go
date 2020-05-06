// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package tuple

import (
	"math/rand"
	"sort"
	"strings"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/util/intmath"
)

type Tuple []string

type Tuples []Tuple

func New(t ...string) Tuple {
	return t
}

func (t Tuple) Get(i int) string {
	return t[i]
}

func (t Tuple) Len() int {
	return len(t)
}

func (t Tuple) String() string {
	return strings.Join(t, "\t")
}

func (t Tuple) Equals(other Tuple) bool {
	for k := 0; k < t.Len(); k++ {
		if t[k] != other[k] {
			return false
		}
	}
	return true
}

func (t Tuple) Sort() {
	sort.Strings(t)
}

func NewTuples() Tuples {
	return make(Tuples, 0)
}

func (ts Tuples) Len() int {
	return len(ts)
}

func (ts Tuples) Sort() {
	sort.Slice(ts, func(i, j int) bool {
		for k := 0; k < ts[i].Len(); k++ {
			if ts[i][k] != ts[j][k] {
				return ts[i][k] < ts[j][k]
			}
		}
		return false
	})
}

func (ts Tuples) Shuffle() {
	for i := 0; i < ts.Len(); i++ {
		j := rand.Intn(i + 1)
		ts[i], ts[j] = ts[j], ts[i]
	}
}

func (ts Tuples) Split(n int) []Tuples {
	ts.Shuffle()
	foldData := make([]Tuples, n)
	cnt := ts.Len()
	k := intmath.Ceil(cnt, n)
	for j := 0; j < n; j++ {
		for i := j * k; i < intmath.Min(cnt, (j+1)*k); i++ {
			foldData[j] = append(foldData[j], ts[i])
		}
	}
	return foldData
}
