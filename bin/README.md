# bin

This directory contains the Caffe2 binary for the medical NER model:
- ner.c2

The model is trained on all [processed annotated data](../data/ner/processed_medical_ner.tsv).

## Model Quality

The F1 score of a NER model trained on the [training data](../data/ner/train_processed_medical_ner.tsv) 
and tested on the [test data](../data/ner/test_processed_medical_ner.tsv) is 0.891. 
For the same test data, the recall at precision is:


| Label                 | R@P 0.2 | R@P 0.4 | R@P 0.6 | R@P 0.8 | R@P 0.9 |
|-----------------------|---------|---------|---------|---------|---------|
| NoLabel               |   1.000 |   1.000 |   1.000 |   0.982 |   0.945 |
| age                   |   0.992 |   0.992 |   0.992 |   0.992 |   0.992 |
| allergy_name          |   0.958 |   0.942 |   0.889 |   0.815 |   0.664 |
| bmi                   |   1.000 |   1.000 |   1.000 |   1.000 |   1.000 |
| cancer                |   0.987 |   0.987 |   0.970 |   0.898 |   0.734 |
| chronic_disease       |   0.995 |   0.985 |   0.973 |   0.902 |   0.788 |
| clinical_variable     |   0.985 |   0.966 |   0.932 |   0.869 |   0.778 |
| contraception_consent |   1.000 |   0.999 |   0.995 |   0.961 |   0.803 |
| ethnicity             |   0.897 |   0.759 |   0.759 |   0.690 |   0.690 |
| gender                |   0.996 |   0.996 |   0.996 |   0.996 |   0.996 |
| language_fluency      |   1.000 |   1.000 |   1.000 |   0.990 |   0.944 |
| lower_bound           |   0.997 |   0.994 |   0.983 |   0.964 |   0.924 |
| pregnancy             |   1.000 |   0.998 |   0.998 |   0.989 |   0.915 |
| technology_access     |   1.000 |   1.000 |   1.000 |   1.000 |   1.000 |
| treatment             |   0.993 |   0.979 |   0.963 |   0.890 |   0.714 |
| upper_bound           |   0.998 |   0.997 |   0.995 |   0.976 |   0.946 |


## License

bin is Apache 2.0 licensed, as found in the [LICENSE file](../LICENSE).
