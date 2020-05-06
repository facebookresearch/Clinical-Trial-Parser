# script

This directory contains scripts for running various clinical-trial modules:
- [cfg_parse.sh](cfg_parse.sh): Parse eligibility criteria with CFG
- [ie_parse.sh](ie_parse.sh): Parse eligibility criteria with IE
- [aact.sh](aact.sh): Download an AACT DB for clinical trials from ClinicalTrials.gov
- [mesh.sh](mesh.sh): Download MeSH descriptors for grounding
- [ingest.sh](ingest.sh): Ingest clinical trial eligibility criteria from the AACT DB to a csv file
- [train_embeddings.sh](train_embeddings.sh): Ingest clinical trial text and train word embeddings
- [search.sh](search.sh): CLI tool to search concepts from a vocabulary

## License

script is Apache 2.0 licensed, as found in the [LICENSE file](../LICENSE).
