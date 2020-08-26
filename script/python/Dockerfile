FROM ubuntu:16.04

# Install Caffe2 + dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
  build-essential \
  git \
  libgoogle-glog-dev \
  libgtest-dev \
  libiomp-dev \
  libleveldb-dev \
  liblmdb-dev \
  libopencv-dev \
  libopenmpi-dev \
  libsnappy-dev \
  openmpi-bin \
  openmpi-doc \
  python-dev \
  python-pip
RUN pip install --upgrade pip
RUN pip install setuptools wheel
RUN pip install future numpy protobuf typing hypothesis pyyaml
RUN apt-get install -y --no-install-recommends \
      libgflags-dev \
      cmake
RUN git clone https://github.com/pytorch/pytorch.git
WORKDIR pytorch
RUN git submodule update --init --recursive
RUN python setup.py install

COPY ./script/python /usr/local/script/python
COPY ./src /usr/local/src
COPY ./data /usr/local/data

RUN pip install -r /usr/local/script/python/requirements.txt
RUN pytext train > /usr/local/src/resources/config/ner.json

USER python
CMD ["python3"]