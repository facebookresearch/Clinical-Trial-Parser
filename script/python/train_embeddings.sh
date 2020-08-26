#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Use clinical study summaries, descriptions, and eligibility criteria
# to train word embedding vectors.
#
# ./script/train_embeddings.sh

set -eu

INGEST_FILE="/usr/local/data/input/clinical_trial_descriptions.txt"
EMBEDDING_FILE="/usr/local/data/embedding/word_embeddings"


echo "Train embedding vectors..."
export PYTHONPATH="$(pwd)/src"
if ! python /usr/local/src/embedding/train_embeddings.py -i "$INGEST_FILE" -o "$EMBEDDING_FILE"
then
  echo "Embedding training failed."
  rm -f "$EMBEDDING_FILE"
  exit 1
fi

# rm $INGEST_FILE
