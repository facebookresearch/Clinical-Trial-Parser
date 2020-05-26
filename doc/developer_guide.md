# Developer Guide

The clinical trial parser library contains tools that can be used to translate 
clinical trial eligibility criteria. For example, it has scripts for downloading data
and running the CFG and IE parsers. The library does not contain publicly available data except 
for 20 clinical trials, which are used to illustrate the functionality of its modules.

## CFG Parser

### Installation steps:

- Install Go from https://golang.org/dl/
- Set `GOPATH` so that the cloned project is in 
`$GOPATH/src/github.com/facebookresearch/Clinical-Trial-Parser`
- Run `./script/cfg_parse.sh` in the project root directory. 
The script will write the parsed relations to [cfg_parsed_clinical_trials.tsv](../data/output/cfg_parsed_clinical_trials.tsv).
- The program parameters can be changed either by changing 
the command line arguments in [cfg_parse.sh](../script/cfg_parse.sh) or 
config parameters in [cfg.conf](../src/resources/config/cfg.conf).

[cfg_parse.sh](../script/cfg_parse.sh) demonstrates how the CFG parser could be used.
Applications should write their own [driver](../src/cmd/cfg/main.go) module.

### Quality improvements:

CFG does not parse all ordinal and numerical criteria. It may also parse some
criteria incorrectly. Errors may be fixed and new capabilities added by:

- Updating the grammar [production rules](../src/ct/parser/production/criterion.go)
by adding new criteria situations. It is also a good practice to add new test cases 
to [interpreter_test.go](../src/ct/parser/interpreter_test.go).
- Updating existing or adding new variables to [variables.csv](../src/resources/variables/variables.csv)
- Updating existing or adding new units to [units.csv](../src/resources/units/units.csv)

## IE Parser

### Installation steps:

- Install Python from https://www.python.org/
- Install Natural Language Toolkit from https://www.nltk.org/
- Install Go from https://golang.org/dl/
- Install [PyText](https://pytext.readthedocs.io/en/master/index.html), which can be done using Anaconda3:
  - Install Anaconda3 from https://docs.anaconda.com/anaconda/install/mac-os/
  - Install PyText with `pip install pytext-nlp`
  - ONNX and Torch may need to be upgraded with  `conda install onnx -c conda-forge` 
  and `conda install pytorch torchvision -c pytorch`
  - Note that PyText has an [issue](https://github.com/facebookresearch/pytext/issues/1365), 
  which affects some users
- Unzip word_embeddings.vec.gz in [data/embedding](../data/embedding)
- Download the MeSH vocabulary using [mesh.sh](../script/mesh.sh)
- Run `./script/ie_parse.sh` in the project root directory. The script will write medical terms and matched concepts
to [ie_parsed_clinical_trials.tsv](../data/output/ie_parsed_clinical_trials.tsv).

The library includes a pre-trained [NER binary](../bin). Drivers and config files are provided 
for illustrative purposes in [src/cmd](../src/cmd) and [src/resources/config](../src/resources/config).
Applications may write their own driver modules.

### Quality improvements:

- The NER model can be improved by adding new training samples
- The NEL module can be improved by
  - A better processing of the extracted NER terms
  - Incorporating a vocabulary that has a high match rate with the eligibility criteria terms
  - Adding synonyms to concepts or new synonyms to the [custom MeSH files](../data/mesh)
  - Implementing term clustering to increase the NEL recall
- Implement RE with negation extraction

## Data

The library includes example scripts [aact.sh](../script/aact.sh) and [ingest.sh](../script/ingest.sh)
for downloading and ingesting clinical trials. While the scripts are provided for convenience, 
applications will most likely need to change them or use other means to do the same. For example, 
`ingest.sh` only samples few trials. An obvious place to start is to change the 'where' clauses.
