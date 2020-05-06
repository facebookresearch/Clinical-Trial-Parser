#!/usr/bin/env python3

# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

import re
import nltk


# number or number range
def is_number(s):
    try:
        float(s.replace(",", "").replace("-", ""))
        return True
    except ValueError:
        return False


class Normalizer:
    def __init__(self):
        self.less_pattern = re.compile(r"[\uFE64\uFF1C\u003C\u227A]")
        self.greater_pattern = re.compile(r"[\uFE65\uFF1E\u003E\u227B]")
        self.less_eq_pattern = re.compile(r"(<\s*/?\s*=|=\s*/?\s*<)")
        self.greater_eq_pattern = re.compile(r"(>\s*/?\s*=|=\s*/?\s*>)")
        self.comparison_pattern = re.compile(r"([=<>\u2264\u2265])")
        self.space_pattern = re.compile(r"\s+")

    def normalize(self, text):
        text = text.strip('"')

        text = re.sub(self.less_pattern, "<", text)
        text = re.sub(self.greater_pattern, ">", text)

        text = text.replace(r"\u2266", r"\u2264")
        text = text.replace(r"\u2267", r"\u2265")

        text = re.sub(self.less_eq_pattern, "\u2264", text)
        text = re.sub(self.greater_eq_pattern, "\u2265", text)

        text = re.sub(self.comparison_pattern, r" \g<1> ", text)

        text = re.sub(self.space_pattern, " ", text)
        text = text.strip().lower()

        return text


class Tokenizer:
    def __init__(self, mask_numbers=True):
        self.mask_numbers = mask_numbers
        self.subtoken_pattern1 = re.compile("(/)")
        self.subtoken_pattern2 = re.compile(r"^(\d+)([a-z].*)")

    def tokenize(self, text):
        tokens = nltk.tokenize.word_tokenize(text)
        tokens = [subtoken for token in tokens for subtoken in re.split(self.subtoken_pattern1, token)]
        tokens = [subtoken for token in tokens for subtoken in re.split(self.subtoken_pattern2, token)]
        tokens = [token for token in tokens if token is not ""]
        if self.mask_numbers:
            tokens = ["@NUMBER" if is_number(token) else token for token in tokens]
        return tokens


class Transformer:
    def __init__(self, mask_numbers=True, download=True):
        self.normalizer = Normalizer()
        self.tokenizer = Tokenizer(mask_numbers=mask_numbers)
        if download:
            nltk.download("punkt")

    def tokenize(self, text):
        text = self.normalizer.normalize(text)
        tokens = self.tokenizer.tokenize(text)
        return tokens

    def transform(self, text):
        tokens = self.tokenize(text)
        return " ".join(tokens)
