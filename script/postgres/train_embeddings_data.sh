#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Use clinical study summaries, descriptions, and eligibility criteria
# to train word embedding vectors.
#
# ./script/train_embeddings.sh

set -eu

INGEST_PATH="/usr/local/data/input"
INGEST_FILE="/usr/local/data/input/clinical_trial_descriptions.txt"
EMBEDDING_FILE="/usr/local/data/embedding/word_embeddings"

mkdir -p "$INGEST_PATH"
touch "$INGEST_FILE"

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

psql -U "$POSTGRES_USER" -d "$DB" -c "$QUERY" > "$INGEST_FILE"

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

psql -U "$POSTGRES_USER" -d "$DB" -c "$QUERY" >> "$INGEST_FILE"

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

psql -U "$POSTGRES_USER" -d "$DB" -c "$QUERY" >> "$INGEST_FILE"

gshuf -o "$INGEST_FILE" < "$INGEST_FILE"
wc -l "$INGEST_FILE"