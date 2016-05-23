#!/bin/sh
cd data
shuf data.csv > data-shuffled.csv
TOTAL_LINES=$(cat data.csv|wc -l)
SPLIT_LINES=$(echo "$TOTAL_LINES/10"|bc)
split -l $SPLIT_LINES --additional-suffix=.csv data-shuffled.csv shuffled-
cat shuffled-aa.csv shuffled-ab.csv shuffled-ac.csv shuffled-ad.csv shuffled-ae.csv shuffled-af.csv shuffled-ag.csv shuffled-ah.csv > shuffled-train.csv
