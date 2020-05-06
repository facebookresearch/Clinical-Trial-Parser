#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Parse clinical-trial eligibility criteria with CFG.
#
# ./script/cfg.sh

set -eu

CMD="src/cmd/cfg/main.go"
CONFIG="src/resources/config/cfg.conf"
INPUT="data/input/clinical_trials.csv"
OUTPUT="data/output/cfg_parsed_clinical_trials.tsv"

if ! go run "$CMD" -conf "$CONFIG" -i "$INPUT" -o "$OUTPUT" -logtostderr
then
  rm -f "$OUTPUT"
  echo "CFG parser failed."
  exit 1
fi
