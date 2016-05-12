#!/bin/sh
mongoexport --db dispotrains --collection stations --out stations.json
mongoexport --db dispotrains --collection statuses --out statuses.json
tar cvf dump.tar stations.json statuses.json
rm stations.json statuses.json
bzip2 -9 dump.tar
mkdir -p /dispotrains/build/app/static/data/
mv dump.tar.bz2 /dispotrains/build/app/static/data/
