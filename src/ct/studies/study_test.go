// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

package studies

import (
	"testing"

	"github.com/facebookresearch/Clinical-Trial-Parser/src/ct/relation"

	"github.com/stretchr/testify/assert"
)

func TestAndCriteriaParse(t *testing.T) {
	a := assert.New(t)

	input := `Inclusion Criteria:

            Male or female, aged 18 to 59 (inclusive).

            NYHA Class of I or II.

            Pre-diabetes: HbA1c >5.7% and BMI ≥ 25 kg/m2.

            Exclusion Criteria:

            Eastern cooperative oncology group is 0-2.`

	expectedInclusions := []relation.Relations{
		relation.Relations{
			relation.Parse(`{"id":"200","name":"age","lower":{"incl":true,"value":"18"},"upper":{"incl":true,"value":"59"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"102","name":"nyha","value":["1", "2"],"variableType":"ordinal"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"25"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":false,"value":"5.7"},"variableType":"numerical"}`),
		},
	}
	expectedExclusions := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["3","4"],"variableType":"ordinal"}`),
	}

	study := NewStudy("ID012345", "Better Health for Everybody", nil, input)
	study.Parse()
	actualInclusionCriteria := study.InclusionCriteria()
	a.Len(actualInclusionCriteria, 4)
	for i, criterion := range actualInclusionCriteria {
		actualInclusions := criterion.Relations()
		actualInclusions.SetScore(0)
		a.Equal(expectedInclusions[i], actualInclusions)
	}

	actualExclusionCriteria := study.ExclusionCriteria()
	a.Len(actualExclusionCriteria, 1)
	actualExclusions := actualExclusionCriteria.Relations()
	actualExclusions.SetScore(0)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestOrCriteriaParse(t *testing.T) {
	a := assert.New(t)

	input := `Inclusion Criteria:

            Pre-diabetes: HbA1c >5.7% or BMI ≥ 25 kg/m2.

            Uncontrolled hypertension, defined as blood pressure (bp) > 140/90 mmhg.

            Eastern Cooperative Oncology Group (ECOG) performance status =< 2.

            Exclusion Criteria:

            Weigh more than 180 pounds.`

	expectedInclusions := []relation.Relations{
		relation.Relations{
			relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"25"},"variableType":"numerical"}`),
			relation.Parse(`{"id":"400","name":"a1c","unit":"%","lower":{"incl":false,"value":"5.7"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"300","name":"sbp","unit":"mmhg","lower":{"incl":false,"value":"140"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"301","name":"dbp","unit":"mmhg","lower":{"incl":false,"value":"90"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"100","name":"ecog","value":["0","1","2"],"variableType":"ordinal"}`),
		},
	}

	expectedExclusions := relation.Relations{
		relation.Parse(`{"id":"202","name":"weight","unit":"lb","upper":{"incl":true,"value":"180"},"variableType":"numerical"}`),
	}

	study := NewStudy("ID012345", "Better Health for Everybody", nil, input)
	study.Parse()
	actualInclusionCriteria := study.InclusionCriteria()
	a.Len(actualInclusionCriteria, 4)
	for i, criterion := range actualInclusionCriteria {
		actualInclusions := criterion.Relations()
		actualInclusions.SetScore(0)
		a.Equal(expectedInclusions[i], actualInclusions)
	}

	actualExclusionCriteria := study.ExclusionCriteria()
	a.Len(actualExclusionCriteria, 1)
	actualExclusions := actualExclusionCriteria[0].Relations()
	actualExclusions.SetScore(0)
	a.Equal(expectedExclusions, actualExclusions)
}

func TestOrdinalNegationCriteriaParse(t *testing.T) {
	a := assert.New(t)

	input := `Inclusion Criteria:

            Height and weight indicating a BMI of at least 18.5 and < 25 kg/m^2 verified at the screening visit.

            Uncontrolled hypertension (defined as blood pressure (bp) > 140/90 mmhg).

            Weigh more than 180 pounds.

            Exclusion Criteria:

            eastern cooperative oncology group (ecog) performance status 2-4.`

	expectedInclusions := []relation.Relations{
		relation.Relations{
			relation.Parse(`{"id":"203","name":"bmi","unit":"kg/m2","lower":{"incl":true,"value":"18.5"},"upper":{"incl":false,"value":"25"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"300","name":"sbp","unit":"mmhg","lower":{"incl":false,"value":"140"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"301","name":"dbp","unit":"mmhg","lower":{"incl":false,"value":"90"},"variableType":"numerical"}`),
		},
		relation.Relations{
			relation.Parse(`{"id":"202","name":"weight","unit":"lb","lower":{"incl":false,"value":"180"},"variableType":"numerical"}`),
		},
	}

	expectedExclusions := relation.Relations{
		relation.Parse(`{"id":"100","name":"ecog","value":["0","1"],"variableType":"ordinal"}`),
	}

	study := NewStudy("ID012345", "Better Health for Everybody", nil, input)
	study.Parse()
	actualInclusionCriteria := study.InclusionCriteria()
	a.Len(actualInclusionCriteria, 4)
	for i, criterion := range actualInclusionCriteria {
		actualInclusions := criterion.Relations()
		actualInclusions.SetScore(0)
		a.Equal(expectedInclusions[i], actualInclusions)
	}

	actualExclusionCriteria := study.ExclusionCriteria()
	a.Len(actualExclusionCriteria, 1)
	actualExclusions := actualExclusionCriteria[0].Relations()
	actualExclusions.SetScore(0)
	a.Equal(expectedExclusions, actualExclusions)
}
