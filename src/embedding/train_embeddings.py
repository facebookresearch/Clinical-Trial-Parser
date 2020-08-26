#!/usr/bin/env python3

# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

"""
Train word embeddings with FastText.

Run: python src/embedding/train_embeddings.py -i <input file> -o <output file>
"""

import os
import fasttext
import numpy as np
from argparse import ArgumentParser

import param as param
from text.transformer import Transformer


def similarity(v1, v2):
    n1 = np.linalg.norm(v1)
    n2 = np.linalg.norm(v2)
    if n1 < param.EPS or n2 < param.EPS:
        return 0
    return np.dot(v1, v2) / (n1 * n2)


def parse_args():
    parser = ArgumentParser()
    parser.add_argument(
        "-i",
        "--input",
        dest="input",
        required=True,
        help="input file for training word embeddings",
        metavar="INPUT",
    )
    parser.add_argument(
        "-o",
        "--output",
        dest="output",
        required=True,
        help="output file for trained word embeddings",
        metavar="OUTPUT",
    )
    return parser.parse_args()


def transform_data(input_file, output_file):
    xfmr = Transformer()
    line_cnt = 0

    with open(output_file, "w", encoding="utf-8") as writer:
        with open(input_file, "r", encoding="utf-8") as reader:
            for line in reader:
                writer.write(xfmr.transform(line) + "\n")
                line_cnt += 1

    print(f"Lines transformed: {line_cnt}")


def train_embeddings(input_file, output_file):
    model = fasttext.train_unsupervised(input_file, **param.options)
    model.save_model(output_file + param.BIN)

    word_writer = open(output_file + param.FREQ, "w", encoding="utf-8")
    vec_writer = open(output_file + param.VEC, "w", encoding="utf-8")

    words, freqs = model.get_words(include_freq=True)
    for w, f in zip(words, freqs):
        word_writer.write(f"{w} {f:d}\n")
        vec = " ".join(format(v, ".6f") for v in model.get_word_vector(w))
        vec_writer.write(f"{w} {vec}\n")

    word_writer.close()
    vec_writer.close()


def test_embeddings(output_file):
    model = fasttext.load_model(output_file + param.BIN)
    for w in param.test_words:
        nearest_neighbors(w, model)
        print()


def nearest_neighbors(w1, model, top=40):
    words, freqs = model.get_words(include_freq=True)
    v1 = model.get_word_vector(w1)

    scores = {}
    for w2, f2 in zip(words, freqs):
        v2 = model.get_word_vector(w2)
        score = similarity(v1, v2)
        if score > 0.5:
            scores[(w2, f2)] = score

    if len(scores) == 0:
        print(f"No nearest neighbors for '{w1}'")
        return

    for k, v in sorted(scores.items(), key=lambda item: item[1], reverse=True)[:top]:
        print(f"{k[0]:20s}  {v:.3f}  {k[1]:7d}")


def main(input_file, output_file):
    tmp_file = output_file + param.TMP
    transform_data(input_file, tmp_file)
    train_embeddings(tmp_file, output_file)
    test_embeddings(output_file)
    os.remove(tmp_file)


if __name__ == "__main__":
    args = parse_args()
    main(args.input, args.output)
