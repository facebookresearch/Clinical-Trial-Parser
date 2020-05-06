#!/usr/bin/env python3

# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

import unittest
from text.transformer import Normalizer
from text.transformer import Transformer


class TestNormalizer(unittest.TestCase):
    norm = Normalizer()

    def test_less(self):
        actual = self.norm.normalize("x <= 4")
        expected = "x ≤ 4"
        self.assertEqual(actual, expected)

    def test_greater(self):
        actual = self.norm.normalize("3>x>/=4 ")
        expected = "3 > x ≥ 4"
        self.assertEqual(actual, expected)


class TestTransformer(unittest.TestCase):
    xfrm = Transformer()

    def test_transform(self):
        actual = self.xfrm.transform("Normal hear,  gen 134dtt4d = /> organ. And/or (marrow): Leukocytes ≥3,000/μL.")
        expected = "normal hear , gen @NUMBER dtt4d ≥ organ . and / or ( marrow ) : leukocytes ≥ @NUMBER / μl ."
        self.assertEqual(actual, expected)

    def test_empty(self):
        actual = self.xfrm.transform("")
        expected = ""
        self.assertEqual(actual, expected)


if __name__ == '__main__':
    unittest.main()
