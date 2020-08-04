#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.
#
# Download a static copy of the AACT database and reconstruct it.
# Set the search path and run an evaluation query for the current year.
# The dump file is deleted at the end.
#
# DATE is YYYYMMDD.
# ./script/aact.sh <DATE>

set -eu

# DATE=$1
DIR="/usr/local/data/aact"
DB=aact
DATE=20200802
ZIP_FILE="${DATE}_clinical_trials.zip"

mkdir -p "$DIR"

# if ! wget -P "$DIR" "https://aact.ctti-clinicaltrials.org/static/static_db_copies/daily/$ZIP_FILE"
# then
#   rm -f "$DIR/$ZIP_FILE"
#   echo "Download failed."
#   exit 1
# fi

# unzip -o -d "$DIR" "$DIR/$ZIP_FILE"

set -e

# psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
# 	GRANT ALL PRIVILEGES ON DATABASE $POSTGRES_DB TO $POSTGRES_USER;
# EOSQL
# dropdb -U "$POSTGRES_USER" --if-exists "$DB"
# createdb -U "$POSTGRES_USER" "$DB"

pg_restore -U "$POSTGRES_USER" -e -v -O -x --dbname="$DB" "${DIR}/postgres_data.dmp"
psql -U "$POSTGRES_USER" -d "$DB" -c "ALTER ROLE $POSTGRES_USER SET search_path TO ctgov,public;"

# Check the download:
YEAR=$(date +'%Y')
echo "Year: $YEAR"

QUERY="
SELECT
    t1.overall_status status,
    COUNT(1) study_count,
    SUM(t1.enrollment) enrollment
FROM (
    SELECT
        nct_id,
        date_part('year', start_date) AS start_year,
        date_part('year', completion_date) AS completion_year,
        enrollment,
        overall_status,
        study_type
    FROM studies
    WHERE
        enrollment < 50000
) t1
JOIN calculated_values t2
    ON t1.nct_id = t2.nct_id
WHERE
    t2.has_us_facility
    AND t1.start_year <= $YEAR
    AND (
        t1.completion_year >= $YEAR
        OR t1.completion_year IS NULL
    )
    AND t1.overall_status = ANY (
        ARRAY[
            'Recruiting',
            'Active, not recruiting',
            'Not yet recruiting',
            'Enrolling by invitation'
        ]
    )
    AND t1.study_type = 'Interventional'
GROUP BY
    t1.overall_status
ORDER BY
    3 DESC;
"

psql -U "$POSTGRES_USER" -d "$DB" -c "$QUERY"

# rm "$DIR/postgres_data.dmp"
