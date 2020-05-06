#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Ingest clinical studies from the aact db to a csv file. 20 sample studies
# are ingested addressing COVID-19 and non-COVID-19 conditions.
#
# ./script/ingest.sh

set -eu

OUTPUT="data/input/clinical_trials.csv"
DB=aact
LIMIT=10

QUERY() {
  echo "\COPY (
  SELECT
      t1.nct_id AS \"#nct_id\",
      t1.brief_title AS title,
      CASE WHEN t2.has_us_facility THEN 'true' ELSE 'false' END AS has_us_facility,
      t3.conditions,
      t4.criteria AS eligibility_criteria
  FROM studies t1
  JOIN calculated_values t2
      ON t1.nct_id = t2.nct_id
  JOIN (
      SELECT
          nct_id,
          STRING_AGG(name, '|' ORDER BY name) AS conditions
      FROM conditions
      GROUP BY
          nct_id
  ) t3
      ON t1.nct_id = t3.nct_id
  JOIN eligibilities t4
      ON t1.nct_id = t4.nct_id
  WHERE
      LOWER(conditions) ${1} '%([^a-z]cov[^a-z]|corona[ v]|covid)%'
      AND t1.study_type = 'Interventional'
      AND t1.overall_status = 'Recruiting'
      AND RANDOM() < 0.2
  ORDER BY
      t1.nct_id DESC
  LIMIT
      ${LIMIT}
  )
  TO STDOUT WITH (FORMAT csv, HEADER)
  "
}

# Extract COVID-19 related trials
psql -U "$USER" -d "$DB" -c "$(QUERY "SIMILAR TO")" > "$OUTPUT"
# Extract non-COVID-19 related trials
psql -U "$USER" -d "$DB" -c "$(QUERY "NOT SIMILAR TO")" >> "$OUTPUT"

wc -l "$OUTPUT"
