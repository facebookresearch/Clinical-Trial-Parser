// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package relation

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/variables"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	a := assert.New(t)

	rs := Relations{
		&Relation{Name: "a", VariableType: variables.Nominal},
		&Relation{Name: "b", VariableType: variables.Ordinal},
	}
	expected := `[{"name":"a","variableType":"nominal","score":0},{"name":"b","variableType":"ordinal","score":0}]`
	actual := rs.JSON()
	a.Equal(expected, actual)
}

func TestNormalize(t *testing.T) {
	a := assert.New(t)

	valueRange := []string{"0", "1", "2", "3", "4"}
	actual := Relation{Name: "a", Lower: &Limit{Incl: false, Value: "0"}, Upper: &Limit{Incl: true, Value: "2"}, Value: []string{"5"}, VariableType: variables.Ordinal}
	expected := Relation{Name: "a", Value: []string{"1", "2"}, VariableType: variables.Ordinal}
	actual.Normalize(valueRange)
	a.Equal(expected, actual)
}

func TestNegate(t *testing.T) {
	a := assert.New(t)

	valueRange := []string{"0", "1", "2", "3", "4"}
	actual := Relation{Name: "a", Value: []string{"3", "4"}, VariableType: variables.Ordinal}
	expected := Relation{Name: "a", Value: []string{"0", "1", "2"}, VariableType: variables.Ordinal}
	actual.Normalize(valueRange)
	actual.Negate(valueRange)
	a.Equal(expected, actual)
}

func TestTransformGoodValue(t *testing.T) {
	a := assert.New(t)

	actual := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "1,000,000"}, VariableType: variables.Numerical}
	expected := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "1000000"}, VariableType: variables.Numerical}
	actual.Transform()
	a.Equal(expected, actual)
}

func TestTransformBadValue(t *testing.T) {
	a := assert.New(t)

	actual := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "xyz"}, VariableType: variables.Numerical, Score: 1}
	expected := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "xyz"}, VariableType: variables.Numerical, Score: 0}
	actual.Transform()
	a.Equal(expected, actual)
}

func TestTransformMissingZero(t *testing.T) {
	a := assert.New(t)

	actual := Relation{ID: variables.Zero, Lower: &Limit{Incl: true, Value: "100,00"}, VariableType: variables.Numerical}
	expected := Relation{ID: variables.Zero, Lower: &Limit{Incl: true, Value: "100000"}, VariableType: variables.Numerical}
	actual.Transform()
	a.Equal(expected, actual)
}

func TestTransformRadix(t *testing.T) {
	a := assert.New(t)

	actual := Relation{ID: variables.Zero, Lower: &Limit{Incl: true, Value: "27,1"}, VariableType: variables.Numerical}
	expected := Relation{ID: variables.Zero, Lower: &Limit{Incl: true, Value: "27.1"}, VariableType: variables.Numerical}
	actual.Transform()
	a.Equal(expected, actual)
}

func TestTransformScientificFormValue(t *testing.T) {
	a := assert.New(t)

	actual := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "1.0  x10^6"}, VariableType: variables.Numerical}
	expected := Relation{ID: variables.Zero, Upper: &Limit{Incl: true, Value: "1.0e6"}, VariableType: variables.Numerical}
	actual.Transform()
	a.Equal(expected, actual)
}
