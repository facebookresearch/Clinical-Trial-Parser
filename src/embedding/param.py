#!/usr/bin/env python3

# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

# fastText options for training word embeddings
options = {
    "model": "cbow",
    "loss": "ns",
    "dim": 100,
    "ws": 5,
    "neg": 5,
    "minCount": 5,
    "bucket": 2000000,
    "t": 1e-5,
    "minn": 0,
    "maxn": 0,
    "lr": 0.05,
    "epoch": 20,
    "thread": 12,
    "verbose": 2,
}

# List of words whose nearest neighbors are computed to inspect the embedding quality
test_words = ["covid-19", "a1c", "cardiomyopathy", "obese", "hemiplegia", "tp53", "cd137", "<", "@NUMBER"]

# Output file extensions
FREQ = ".freq"
VEC = ".vec"
BIN = ".bin"
TMP = ".tmp"

# Vector is approximated by zero if its norm is less than EPS
EPS = 1.0e-6
