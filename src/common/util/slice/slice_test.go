// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package slice

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/common/col/set"

	"github.com/stretchr/testify/assert"
)

func TestSetToSlice(t *testing.T) {
	a := assert.New(t)

	input := set.New("21", "3", "a", "12")
	expected := []string{"12", "21", "3", "a"}
	actual := SetToSlice(input)
	a.Equal(expected, actual)
}

func TestNumberSetToSlice(t *testing.T) {
	a := assert.New(t)

	input := set.New("21", "3", "5", "12")
	expected := []string{"3", "5", "12", "21"}
	actual := SetToSlice(input)
	a.Equal(expected, actual)
}
