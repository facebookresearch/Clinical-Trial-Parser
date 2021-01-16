#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Download the MeSH descriptors to data/mesh. The first argument is
# the optional latest production year with the default value "2021".
#
# ./script/mesh.sh [<year>]

set -eu

PRODUCTION_YEAR=${1:-"2021"}
DESCRIPTOR=desc${PRODUCTION_YEAR}.xml

if ! curl ftp://nlmpubs.nlm.nih.gov/online/mesh/MESH_FILES/xmlmesh/"$DESCRIPTOR" -o data/mesh/descriptor.xml
then
  echo "MeSH descriptor download failed; the latest production year may be old: $PRODUCTION_YEAR"
  exit 1
fi