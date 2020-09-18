import pathlib
import re
import argparse
from typing import NamedTuple, List, Iterator, Dict, Tuple
import pickle



def transform_data(input_data_path, output_data_path):

    input_data_path = pathlib.Path(input_data_path)
    output_data_path = pathlib.Path(output_data_path)

    entities = r"(\d{1,3}):(\d{1,3}):([a-z]*_?[a-z]*)"
    negated_text = r"\d{1,3}:\d{1,3}:[a-z]*_?[a-z]*,?\s?"

    with input_data_path.open() as file:
        spacy_list = []
        for n,line in enumerate(file):
            matched_entities = re.findall(entities, line)
            
            matched_entities_formatted = []
            for matched_entity in matched_entities:
                matched_entity_0 = int(matched_entity[0])-1
                matched_entity_1 = int(matched_entity[1])-1
                matched_entity_2 = str(matched_entity[2])
                matched_entities_formatted.append((matched_entity_0, matched_entity_1, matched_entity_2))

            matched_text = re.sub(negated_text, '', line).replace('\n', '')
            spacy_tuple = (matched_text, {"entities": matched_entities_formatted})
            spacy_list.append(spacy_tuple)

    with open(output_data_path, 'wb') as f:
        pickle.dump(spacy_list, f)


if __name__ == "__main__":
    parser = argparse.ArgumentParser()
    parser.add_argument(
            '--input_data_path',
            help="Path to the input data."
    )

    parser.add_argument(
            '--output_data_path',
            help="Path to the output data."
    )

    args = parser.parse_args()
    transform_data(
        args.input_data_path,
        args.output_data_path)