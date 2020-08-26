#!/usr/bin/env bash
# apt install build-essential \
#     && apt-get install manpages-dev \
#     && pip install pytext-nlp \
#     && pytext train > /usr/local/src/resources/config/ner.json

# && conda install onnx -c conda-forge
# pip install pytext-nlp==0.3.2 \
    # && conda install onnx -c conda-forge -y \
    # && conda install pytorch torchvision -c pytorch -y \
    # && pip install -r /usr/local/script/python/requirements.txt \
    # && pytext train > /usr/local/src/resources/config/ner.json

pytext train > /usr/local/src/resources/config/ner.json