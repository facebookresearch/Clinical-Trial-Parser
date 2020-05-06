// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package parser

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"

	"github.com/stretchr/testify/assert"
)

func TestOneNumericalVariableInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "a1c greater than or equal to 5.0%."
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":true,"value":"5.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestOneNumericalVariableRangeInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "screening visit a1c ≥ 7.0% and ≤ 12.5%"
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":true,"value":"7.0"},"upper":{"incl":true,"value":"12.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestOneCategoricalVariableRangeInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "eastern cooperative oncology group (ecog) performance status  0-3"
	expected := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2","3"],"variableType":"ordinal"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestOneCategoricalVariableSetInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "eastern cooperative oncology group (ecog) performance status of 0, 1, or 2."
	expected := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2"],"variableType":"ordinal"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestNumberUnitInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "absolute neutrophil count of >= 1000 mm3"
	expected := relation.Relations{
		relation.Parse(`{"id":"408","name":"anc","unit":"cells/ul","lower":{"incl":true,"value":"1000"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestReverseOrderInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "patient < 30.0 kilograms (kg) in weight"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","upper":{"incl":false,"value":"30.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestMisspelledVariableNameInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "weigh less than 110 pounds"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"lb","upper":{"incl":false,"value":"110"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestMissingVariableNameInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "patients less than 31 pounds"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"lb","upper":{"incl":false,"value":"31"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoVariablesInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "weight > 50 kg and a1c <= 9.5%"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","upper":{"incl":true,"value":"9.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoVariablesWithExtraWordsInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "pre-diabetes: hba1c >5.7% and bmi ≥ 25 kg/ m2"
	expected := relation.Relations{
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"25"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":false,"value":"5.7"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoVariablesWithMoreStuffInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "weight ≥ 50 kg; with a body mass index ≤ 45 kg/m2"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":true,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","upper":{"incl":true,"value":"45"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestRangeWithBetweenAndInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "a body mass index between 18.0kg /m2 and 31.0 kg/m2, and a total body weight >50 kg"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.0"},"upper":{"incl":true,"value":"31.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestRangeWithBetweenOrInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "a body mass index between 18.0kg / m2 and 31.0, or a total body weight >50 kg"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.0"},"upper":{"incl":true,"value":"31.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualOrRels.Process()
	actualOrRels.SetScore(0)

	a.Empty(actualAndRels)
	a.Equal(expected, actualOrRels)
}

func TestTwoVariablesWithAndOrInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "a body mass index between 18.0kg/m2 and 31.0 kg/m2, and/or a total body weight >50 kg"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.0"},"upper":{"incl":true,"value":"31.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualOrRels.Process()
	actualOrRels.SetScore(0)

	a.Empty(actualAndRels)
	a.Equal(expected, actualOrRels)
}

func TestRangeWithDashInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "weight >50 kg; with a body mass index 18.0 - 31.0 kg/m2"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.0"},"upper":{"incl":true,"value":"31.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestRangeWithToInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "a body mass index 18.0kg/m2 to 31.0 kg/m2"
	expected := relation.Relations{
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.0"},"upper":{"incl":true,"value":"31.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestRangeWithOrInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "	must have a body mass index (bmi) of greater than (>) 18.5 or less than equal to (<=) 35 kilogram per meter square (kg/m^2)"
	expected := relation.Relations{
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":false,"value":"18.5"},"upper":{"incl":true,"value":"35"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestParenthesisInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "eastern cooperative oncology group (ecog) performance status =< 2 (karnofsky >= 50%)"
	expected := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2"],"variableType":"ordinal"}`),
		relation.Parse(`{"id":"600","name":"karnofsky_score","unit":"%","lower":{"incl":true,"value":"50"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoVariablesInParenthesisInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "uncontrolled hypertension (systolic blood pressure (sbp) >140mmhg, diastolic blood pressure (dbp) >90mmhg)"
	expected := relation.Relations{
		relation.Parse(`{"id":"300","name":"sbp","unit":"mmhg","lower":{"incl":false,"value":"140"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"301","name":"dbp","unit":"mmhg","lower":{"incl":false,"value":"90"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestOneVariableWithTwoValuesInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "uncontrolled hypertension, defined as blood pressure (bp) > 140/90"
	expected := relation.Relations{
		relation.Parse(`{"id":"300","name":"sbp","lower":{"incl":false,"value":"140"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"301","name":"dbp","lower":{"incl":false,"value":"90"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestVariableWithAcronymUnitInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "aspartate aminotransferase (ast) =< 2.5 x uln (to be performed within14 days prior to day 1 of protocol therapy unless otherwise stated)"
	expected := relation.Relations{
		relation.Parse(`{"id":"411","name":"ast","unit":"uln","upper":{"incl":true,"value":"2.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestCompositeVariableWithSpaceInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "aspartate aminotransferase (ast)/alanine aminotransferase (alt) ≤ 2.0 x upper limits of normal"
	expected := relation.Relations{
		relation.Parse(`{"id":"411","name":"ast","unit":"uln","upper":{"incl":true,"value":"2.0"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"412","name":"alt","unit":"uln","upper":{"incl":true,"value":"2.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestCompositeVariableWithoutSpaceInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "sgot/sgpt ≤ 3 x laboratory normal or ≤ 5 x laboratory normal if something else"
	expected := relation.Relations{
		relation.Parse(`{"id":"411","name":"ast","unit":"uln","upper":{"incl":true,"value":"3"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"412","name":"alt","unit":"uln","upper":{"incl":true,"value":"3"},"variableType":"numerical"}`),
	}

	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestCompositeVariableWithOrInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "aspartate aminotransferase (ast) or alanine aminotransferase (alt) ≤ 2.0 x upper limits of normal"
	expected := relation.Relations{
		relation.Parse(`{"id":"411","name":"ast","unit":"uln","upper":{"incl":true,"value":"2.0"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"412","name":"alt","unit":"uln","upper":{"incl":true,"value":"2.0"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestRatioInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "sgot/sgpt ratio > 2.0"
	expected := relation.Relations{
		relation.Parse(`{"id":"414","name":"ast/alt_ratio","lower":{"incl":false,"value":"2.0"},"variableType":"numerical"}`),
	}

	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestVariableWithMultipleNumbersInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "hba1c of 7.0-9.5% (53-80 mmol/mol) (both inclusive) as assessed by central laboratory"
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":true,"value":"7.0"},"upper":{"incl":true,"value":"9.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestVariableNamesWithSharedWordsShortInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "patients with serum triglyceride levels >300 mg/dl"
	expected := relation.Relations{
		relation.Parse(`{"id":"506","name":"triglyceride_level","unit":"mg/dl","lower":{"incl":false,"value":"300"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestVariableNamesWithSharedWordsLongInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "patients with fasting serum triglyceride levels> 200 mg/dl"
	expected := relation.Relations{
		relation.Parse(`{"id":"505","name":"fasting_triglyceride_level","unit":"mg/dl","lower":{"incl":false,"value":"200"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoOrVariablesInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "weight > 50 kg or a1c <= 9.5%"
	expected := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"kg","lower":{"incl":false,"value":"50"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","upper":{"incl":true,"value":"9.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualOrRels.Process()
	actualOrRels.SetScore(0)

	a.Empty(actualAndRels)
	a.Equal(expected, actualOrRels)
}

func TestMultiNumberlVariableInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "patients without type 2 diabetes or with type 2 diabetes who have a a1c less than or equal to 7.9%"
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","upper":{"incl":true,"value":"7.9"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestComparisonInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "hba1c of 8 or greater"
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","lower":{"incl":true,"value":"8"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestToRangeInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "hba1c >6.0% to <=11.5%"
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":false,"value":"6.0"},"upper":{"incl":true,"value":"11.5"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestUnknownVariableInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "screening a1c ≥ 10% or unknown xyz < 2.0 mg/dl."
	expected := relation.Relations{
		relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":true,"value":"10"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualOrRels.Process()
	actualOrRels.SetScore(0)

	a.Empty(actualAndRels)
	a.Equal(expected, actualOrRels)
}

func TestColonInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "histologically confirmed or suspected diagnosis of 1 of the following: performance status - ecog 0-2"
	expected := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2"],"variableType":"ordinal"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestTwoSameVariablesInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "bmi between 20 to 25 kg/m^2 or bmi 30 to 35 kg/m^2"
	expected := relation.Relations{
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"20"},"upper":{"incl":true,"value":"25"},"variableType":"numerical"}`),
		relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"30"},"upper":{"incl":true,"value":"35"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualOrRels.Process()
	actualOrRels.SetScore(0)

	a.Empty(actualAndRels)
	a.Equal(expected, actualOrRels)
}

func TestRomanNumeralsInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "nyha class of ii, iii or iv."
	expected := relation.Relations{
		relation.Parse(`{"id":"102","name":"nyha","value":["2","3","4"],"variableType":"ordinal"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestSimpleCholesterolInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "normal cholesterol <= 200 mg/dl"
	expected := relation.Relations{
		relation.Parse(`{"id":"500","name":"total_cholesterol","unit":"mg/dl","upper":{"incl":true,"value":"200"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestNumberInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "leukocytes ≥ 100 × 109 / l test"
	expected := relation.Relations{
		relation.Parse(`{"id":"404","name":"wbc","unit":"cells/l","lower":{"incl":true,"value":"100e9"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.Transform()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}

func TestNonCompositeVariableWithSlashInterpreter(t *testing.T) {
	a := assert.New(t)

	input := "pao2 /fio2 < 200."
	expected := relation.Relations{
		relation.Parse(`{"id":"904","name":"pf_ratio","unit":"mmhg","upper":{"incl":false,"value":"200"},"variableType":"numerical"}`),
	}
	actualOrRels, actualAndRels := interpreter.Interpret(input)
	actualAndRels.Process()
	actualAndRels.SetScore(0)

	a.Empty(actualOrRels)
	a.Equal(expected, actualAndRels)
}
