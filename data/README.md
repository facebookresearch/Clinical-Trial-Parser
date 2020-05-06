# data

The data directory contains:

- Annotated word labeling data for training and testing named-entity recognition (NER) models
- Medical word embeddings
- Sample input and output data for clinical trials
- Custom medical concepts and synonyms

## Annotated Word Labeling Data

This annotated medical NER dataset is composed of the eligibility requirements 
from 3.3K randomly sampled interventional trials in the United States, in recruiting status. 
The eligibility requirements were split by line into 50K samples. Each sample was then 
annotated by professional annotators.

In total 120K slots were labeled along the following distribution of labels:

| Label               | Count|
| :---                | ---: |
|treatment            | 30972|
|chronic_disease      | 26212|
|upper_bound          | 13967|
|lower_bound          | 13633|
|clinical_variable    | 13255|
|cancer               |  9344|
|gender               |  3661|
|pregnancy            |  2773|
|age                  |  2616|
|allergy_name         |  1887|
|contraception_consent|  1603|
|language_fluency     |   482|
|bmi                  |   287|
|technology_access    |   132|
|ethnicity            |    82|

Included files are:
- [medical_ner.tsv](ner/medical_ner.tsv): Unprocessed annotated medical NER samples. Slots are grouped by line and NCT ID
- [processed_medical_ner.tsv](ner/medical_ner.tsv): A processed version of the same samples. Text is normalized and numbers tagged as @NUMBER
- [train_processed_medical_ner.csv](ner/train_processed_medical_ner.tsv), [test_processed_medical_ner.csv](ner/test_processed_medical_ner.tsv): Training and test datasets, with an 90-10 split of deduplicated processed_medical_ner.tsv

## Word Embeddings

### Parameters

Word vectors are trained on clinical trial descriptions and eligibility criteria using fastText.
The dimension of the word vectors is 100. The other fastText parameters are:

| Parameter | Value |
|:--- | ---: |
| model | cbow |
|loss | ns |
| dim | 100 |
| ws | 5 |
| neg| 5 |
| minCount | 5 |
| bucket | 2000000 |
| t | 1e-5 |
| minn | 0 |
| maxn | 0 |
| lr| 0.05 |
| epoch | 20 |

### Nearest neighbors of example words

**'covid-19'**:

| Word                | Similarity | Frequency |
|:--- | :---: | ---:|
|covid-19            |  1.000   |  2473|
|sars-cov-2          |  0.920   |   802|
|sars-cov2           |  0.785   |    94|
|coronavirus         |  0.769   |   716|
|2019-ncov           |  0.747   |   209|
|sars                |  0.734   |   341|
|covid19             |  0.732   |    75|
|covid               |  0.705   |   139|
|pandemic            |  0.656   |  1198|
|mers-cov            |  0.649   |    59|
|hantavirus          |  0.630   |    75|
|covid-2019          |  0.625   |    13|
|isaric              |  0.623   |     7|
|cov2                |  0.619   |     7|
|ncov                |  0.619   |    17|
|ards                |  0.611   |  3346|
|hyperinflammation   |  0.598   |    27|
|outbreak            |  0.593   |   679|
|influenza-like      |  0.592   |   234|
|convalescent        |  0.590   |   197|
|covid-19+           |  0.586   |    12|
|influenza           |  0.585   | 14706|
|cov-2               |  0.585   |    21|
|sars-cov            |  0.584   |    66|
|oseltamivir         |  0.580   |   707|
|ifv                 |  0.579   |    17|
|lcri                |  0.575   |     7|
|rsv                 |  0.573   |  1608|
|favipiravir         |  0.572   |    86|
|remdesivir          |  0.569   |    28|
|mers                |  0.563   |    72|
|lassa               |  0.563   |    73|
|quarantine          |  0.560   |   108|
|monkeypox           |  0.547   |    37|
|sari                |  0.544   |    98|
|h1n1                |  0.539   |  1346|
|h1n1v               |  0.535   |    30|
|ntzx                |  0.532   |     7|
|sars-covid-19       |  0.532   |     9|
|sars-cov-1          |  0.531   |     9|

**'A1c'**:

| Word                | Similarity | Frequency |
|:--- | :---: | ---:|
|a1c                 |  1.000   |  4370|
|hba1c               |  0.900   | 14541|
|hgba1c              |  0.839   |   353|
|glycosylated        |  0.837   |  1479|
|glycated            |  0.822   |  1057|
|hga1c               |  0.816   |   170|
|ha1c                |  0.751   |    57|
|hb1ac               |  0.725   |    51|
|hemoglobin          |  0.692   | 31993|
|7.5-9.0             |  0.692   |     6|
|hbalc               |  0.688   |    52|
|hbac1               |  0.673   |    12|

**'<'**:

| Word                | Similarity | Frequency |
|:--- | :---: | ---:|
|<                 |    1.000 |  215427|
|>                 |    0.951 |  284378|
|/                 |    0.869 | 1025687|
|≥                 |    0.869 |  243715|
|dl                |    0.795 |   82963|
|≤                 |    0.784 |  126941|
|;                 |    0.763 |  545143|
|@NUMBER           |    0.761 | 6072529|
|l                 |    0.756 |   66810|
|-                 |    0.744 | 3713603|
|creatinine        |    0.742 |   69367|
|min               |    0.727 |   49914|
|greater           |    0.719 |  110059|
|less              |    0.703 |  123778|

## Input and Output Data

- [Input](input/clinical_trials.csv) is a sample of 20 recent clinical trials, half of which are for COVID-19 conditions
- [Output](output) contains parsed eligibility criteria for the sampled clinical trials

## Custom medical concepts and synonyms

MeSH is augmented with custom concepts and synonyms to improve eligibility criteria parsing. 
For example, new vocabulary concepts are created for the simplest domains and synonyms are added
to existing MeSH concepts.

## Acknowledgement

Thanks to the [Clinical Trials Transformation Initiative](https://www.ctti-clinicaltrials.org/) (CTTI) 
for providing the [Aggregate Analysis of ClinicalTrials.gov](https://aact.ctti-clinicaltrials.org/) (AACT) 
Database for the registered clinical studies at [ClinicalTrials.gov](https://clinicaltrials.gov/ct2/home). 
The sample trials were ingested from the AACT database using the daily static DB copy of 2020-04-16. 
These [sample studies](input/clinical_trials.csv) include a subset of trial information as shown on 
[ClinicalTrials.gov](https://clinicaltrials.gov/ct2/home), for illustrative purposes only. 
The [architecture description](../doc/architecture.md) provides a detailed explanation of 
the modifications made in the sample [output directory](output).

## License

data is Apache 2.0 licensed, as found in the [LICENSE file](../LICENSE).
