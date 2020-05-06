#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Download the MeSH descriptors to data/mesh.
#
# ./script/mesh.sh

set -eu

DESCRIPTOR=desc2020.xml

curl ftp://nlmpubs.nlm.nih.gov/online/mesh/MESH_FILES/xmlmesh/${DESCRIPTOR} -o data/mesh/${DESCRIPTOR}
