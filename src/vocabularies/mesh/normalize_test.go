// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package mesh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeHepatitis(t *testing.T) {
	a := assert.New(t)

	input := " hepatitis c evidenced"
	expected := "c hepatitis"
	actual, _ := Normalize(input)

	a.Equal(expected, actual)
}

func TestNormalizeDBP(t *testing.T) {
	a := assert.New(t)

	input := "diastolic blood pressure [dbp] test"
	expected := "blood diastolic pressure test"
	actual, _ := Normalize(input)

	a.Equal(expected, actual)
}

func TestNormalizeFilter(t *testing.T) {
	a := assert.New(t)

	input := "w1 and/or w2"
	expected := "w1 w2"
	_, actual := Normalize(input)

	a.Equal(expected, actual)
}
