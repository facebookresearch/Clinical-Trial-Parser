#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# For a given term or phrase, search for matching concepts from a vocabulary.
# Search terms are read from console.
#
# ./script/search.sh

set -eu

CMD="src/cmd/search/main.go"
CONFIG="src/resources/config/search.conf"

go run "$CMD" -conf "$CONFIG" -logtostderr
