#!/usr/bin/env python3

# python src/ie/ner.py -m <spacy model name (self trained: ner_clinical_trial_md)> -i <input file> -o <output file>


import json
from argparse import ArgumentParser
import spacy
import pathlib



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



def main(model_name, input_file, output_file):

    nlp = spacy.load("en_ner_clinicaltrialgov_lg")

    input_data_path = pathlib.Path(input_file)
    output_data_path = pathlib.Path(output_file)


    with open(input_data_path, "r") as reader:
        with open(output_data_path, "w") as writer:
            line = reader.readline().strip()
            writer.write(line + "\tdetected_slots\n")  # header
            text_list = []
            for line in reader:
                fields = line.strip().split("\t")
                if len(fields) != 3:
                    print(f"bad row: {fields}")
                    continue
                text_list.append(line)
                doc = nlp(fields[2])
                for ent in doc.ents:
                    ent_dict = {"label": f"word_scores:{ent.label_}", "term": ent.text, "score": 1.0}
                    writer.write("\t".join(fields + [json.dumps(ent_dict)]) + "\n")


    # with input_data_path.open() as file:
    #     spacy_list = []
    #     for n,line in enumerate(file):
    #         spacy_list.append(line)
    # print(spacy_list[:5])

    # nlp = spacy.load("en_ner_clinicaltrialgov_lg")
    # docs = list(nlp.pipe())
    # doc = nlp("Patients with respiratory distress syndrome with arterial partial pressure of oxygen / fraction of inspired oxygen (PaO2 / FiO2) <200 mm/Hg or")
    # print(f'tokens: {[token.text for token in doc]}')
    # print(f'coarse-grained part of speech: {[token.pos_ for token in doc]}')
    # print(f'fine-grained part of speech: {[token.tag_ for token in doc]}')
    # print(f'dependencies: {[token.dep_ for token in doc]}')
    # print(f'entities: {[(ent.text, ent.start_char, ent.end_char, ent.label_) for ent in doc.ents]}')



if __name__ == "__main__":
    args = parse_args()
    main(args.model, args.input, args.output)
