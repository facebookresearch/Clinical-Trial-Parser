#!/usr/bin/env python3

# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

"""
Process annotated eligibility criteria to a format suitable for training
and testing NER models with PyText. Duplicate samples are removed.

Run: python src/ie/process_ner_data.py -i <input file> -o <output file>

For example:
export PYTHONPATH="$(pwd)/src"
python src/ie/process_ner_data.py -i data/ner/medical_ner.tsv -t data/ner/train_processed_medical_ner.tsv -v data/ner/test_processed_medical_ner.tsv
"""

import random
import os.path
from argparse import ArgumentParser
from text.transformer import Transformer

# Text transformer for eligibility criteria
xfmr = Transformer()

# Ratio of samples assigned to the test set
test_ratio = 0.1


def parse_args():
    parser = ArgumentParser()
    parser.add_argument(
        "-i",
        "--input",
        dest="input",
        required=True,
        help="input file",
        metavar="INPUT",
    )
    parser.add_argument(
        "-t",
        "--train",
        dest="train",
        required=True,
        help="output train file",
        metavar="TRAIN",
    )
    parser.add_argument(
        "-v",
        "--test",
        dest="test",
        required=True,
        help="output test file",
        metavar="TEST",
    )
    return parser.parse_args()


def transform(slots, text):
    new_labels = []
    new_text = []
    previous_hi = 0
    for slot in slots:
        values = slot.split(":")
        lo = int(values[0]) - 1
        hi = int(values[1]) - 1
        label = values[2]
        if lo > previous_hi:
            new_text.append(xfmr.transform(text[previous_hi:lo]))
            new_labels.append("")
        new_text.append(xfmr.transform(text[lo:hi]))
        new_labels.append(label)
        previous_hi = hi
    if previous_hi < len(text):
        new_text.append(xfmr.transform(text[previous_hi:]))
        new_labels.append("")

    new_text = [v.strip() for v in new_text]

    text = ""
    labels = ""
    for i, v in enumerate(new_text):
        label = new_labels[i]
        if v == "":
            if label == "":
                continue
            else:
                print("bad transform:")
                print("original: " + labels + "\t" + text)
                print("next text: " + " ".join(new_text))
                print("new labels: " + " ".join(new_labels))
                exit(1)
        if i > 0:
            text += " "
        if label == "":
            text += v
        else:
            lo = len(text) + 1
            text += v
            hi = len(text) + 1
            if labels != "":
                labels += ","
            labels += f"{lo}:{hi}:{label}"

    return labels + "\t" + text


def main(input_file, train_file, test_file):
    samples = set()
    trials = set()
    histo = dict()
    line_cnt = 0
    slot_cnt = 0
    train_cnt = 0
    test_cnt = 0
    with open(input_file, "r", encoding="utf-8") as reader:
        with open(train_file, "w", encoding="utf-8") as train_writer:
            with open(test_file, "w", encoding="utf-8") as test_writer:
                for line in reader:
                    line_cnt += 1
                    values = line.strip().split("\t")
                    sample = values[1] + "\t" + values[2]
                    trials.add(values[0])
                    if sample not in samples:
                        samples.add(sample)
                        slots = values[1].split(",")
                        slot_cnt += len(slots)
                        new_sample = transform(slots, values[2])
                        for slot in slots:
                            label = slot.split(":")[2]
                            histo[label] = histo.get(label, 0) + 1
                        if random.random() < test_ratio:
                            test_writer.write(new_sample + "\n")
                            test_cnt += 1
                        else:
                            train_writer.write(new_sample + "\n")
                            train_cnt += 1

    print(f"Train count: {train_cnt}, test count: {test_cnt} ({100 * test_cnt / len(samples):.1f}%)")
    print(f"Lines read: {line_cnt}, slots: {slot_cnt}, samples: {len(samples)}, trials: {len(trials)}")
    print()

    items = list(histo.items())
    items.sort(key=lambda i: i[1], reverse=True)
    print("label".ljust(22, " ") + "count")
    tot = 0
    for k, v in items:
        print(f"{k:21s} {v:5d}")
        tot += v
    print("total".ljust(21, " ") + f"{tot:6d}")


if __name__ == "__main__":
    args = parse_args()
    main(args.input, args.train, args.test)
