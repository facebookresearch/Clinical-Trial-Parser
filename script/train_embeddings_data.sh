#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Use clinical study summaries, descriptions, and eligibility criteria
# to train word embedding vectors.
#
# ./script/train_embeddings.sh

set -eu

INGEST_FILE="data/input/clinical_trial_descriptions.txt"
EMBEDDING_FILE="data/embedding/word_embeddings"

DB=aact
LENGTH_LIMIT=20

echo "Ingest brief trial descriptions..."
QUERY="\COPY (
SELECT
    TRIM(REGEXP_REPLACE(description, '\s+', ' ', 'g'))
FROM
    brief_summaries
WHERE
    LENGTH(description) > $LENGTH_LIMIT
)
TO STDOUT"

psql -U "$USER" -d "$DB" -c "$QUERY" > "$INGEST_FILE"

echo "Ingest detailed trial descriptions..."
QUERY="\COPY (
SELECT
    TRIM(REGEXP_REPLACE(description, '\s+', ' ', 'g'))
FROM
    detailed_descriptions
WHERE
    LENGTH(description) > $LENGTH_LIMIT
)
TO STDOUT"

psql -U "$USER" -d "$DB" -c "$QUERY" >> "$INGEST_FILE"

echo "Ingest trial eligibility criteria..."
QUERY="\COPY (
SELECT
    TRIM(REGEXP_REPLACE(criteria, '\s+', ' ', 'g'))
FROM
    eligibilities
WHERE
    LENGTH(criteria) > $LENGTH_LIMIT
)
TO STDOUT"

psql -U "$USER" -d "$DB" -c "$QUERY" >> "$INGEST_FILE"

gshuf -o "$INGEST_FILE" < "$INGEST_FILE"
wc -l "$INGEST_FILE"

# echo "Train embedding vectors..."
# export PYTHONPATH="$(pwd)/src"
# if ! python src/embedding/train_embeddings.py -i "$INGEST_FILE" -o "$EMBEDDING_FILE"
# then
#   echo "Embedding training failed."
#   rm -f "$EMBEDDING_FILE"
#   exit 1
# fi

# rm $INGEST_FILE
