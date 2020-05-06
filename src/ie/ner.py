#!/usr/bin/env python3
#
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Extract medical terms from clinical-trial eligibility criteria with a NER model.
#
# python src/ie/ner.py -m <model file> -i <input file> -o <output file>


import json
from itertools import groupby
from argparse import ArgumentParser
import numpy as np

from caffe2.python import workspace
from caffe2.python.predictor import predictor_exporter
from text.transformer import Transformer

xfmr = Transformer()


def parse_args():
    parser = ArgumentParser()
    parser.add_argument(
        "-m",
        "--model",
        dest="model",
        required=True,
        help="model bin for NER",
        metavar="MODEL",
    )
    parser.add_argument(
        "-i",
        "--input",
        dest="input",
        required=True,
        help="input file for embeddings",
        metavar="INPUT",
    )
    parser.add_argument(
        "-o",
        "--output",
        dest="output",
        required=True,
        help="output file for embeddings",
        metavar="OUTPUT",
    )
    return parser.parse_args()


def tokenize(text, add_bos=False):
    tokens = xfmr.tokenize(text)
    if add_bos:
        tokens = ["__BEGIN_OF_SENTENCE__"] + tokens
    return tokens


class OfflinePredictor:
    def __init__(self, model_file):
        self.predict_net = predictor_exporter.prepare_prediction_net(
            model_file, db_type="minidb"
        )

    def predict_pytext(self, tokens):
        tokens_vals = np.array([tokens], dtype=np.str)
        tokens_lens = np.array([len(tokens)], dtype=np.int64)
        workspace.FeedBlob("tokens_vals_str:value", tokens_vals)
        workspace.FeedBlob("tokens_lens", tokens_lens)

        workspace.RunNet(self.predict_net)
        prediction = {
            str(blob): workspace.blobs[blob]
            for blob in self.predict_net.external_outputs
        }
        return prediction


def predictions_by_word(predictions, split_prediction_text):
    word_prediction_array = []
    for index, word in enumerate(split_prediction_text, start=0):
        predictions_by_label = [
            {"score": value_array[0][index][0], "label": label}
            for label, value_array in predictions.items()
        ]
        exps_sum = sum(
            [np.exp(prediction["score"]) for prediction in predictions_by_label]
        )
        soft_maxed_predictions_by_label = [
            {
                "score": np.exp(prediction["score"]) / exps_sum,
                "label": prediction["label"],
            }
            for prediction in predictions_by_label
        ]
        best_prediction = max(soft_maxed_predictions_by_label, key=lambda x: x["score"])
        word_prediction_array.append(
            {
                "word": word,
                "label": best_prediction["label"],
                "score": best_prediction["score"],
            }
        )
    return word_prediction_array


def group_slots(word_prediction_array):
    prev_prediction = word_prediction_array[0]
    lumped_predictions = [[prev_prediction]]
    del word_prediction_array[0]
    for word_prediction in word_prediction_array:
        if word_prediction["label"] == prev_prediction["label"]:
            prediction_set = lumped_predictions.pop()
        else:
            prediction_set = []
        prediction_set.append(word_prediction)
        lumped_predictions.append(prediction_set)
        prev_prediction = word_prediction

    slots = []
    for prediction_set in lumped_predictions:
        slot_text = " ".join(map(lambda x: x["word"], prediction_set))
        slot_score = np.mean(list(map(lambda x: float(x["score"]), prediction_set)))
        slot_label = prediction_set[0]["label"]
        slots.append({"text": slot_text, "score": slot_score, "label": slot_label})

    grouped_slots = {}
    for label, slot in groupby(slots, lambda x: x["label"]):
        if label not in grouped_slots:
            grouped_slots[label] = []
        slot = list(slot)[0]
        grouped_slots[label].append([slot["score"], slot["text"]])
    return grouped_slots


def predict(offline_predictor, prediction_text):
    tokens = tokenize(prediction_text)
    predictions = offline_predictor.predict_pytext(tokens)
    word_prediction_array = predictions_by_word(predictions, tokens)
    return group_slots(word_prediction_array)


def main(model_file, input_file, output_file):
    offline_predictor = OfflinePredictor(model_file)
    with open(input_file, "r") as reader:
        with open(output_file, "w") as writer:
            line = reader.readline().strip()
            writer.write(line + "\tdetected_slots\n")  # header
            for line in reader:
                fields = line.strip().split("\t")
                if len(fields) != 3:
                    print(f"bad row: {fields}")
                    continue
                grouped_slots = predict(offline_predictor, fields[2])
                filtered_grouped_slots = {
                    slot_name: slots
                    for slot_name, slots in grouped_slots.items()
                    if slot_name != "word_scores:NoLabel"
                }
                writer.write("\t".join(fields + [json.dumps(filtered_grouped_slots)]) + "\n")


if __name__ == "__main__":
    args = parse_args()
    main(args.model, args.input, args.output)
