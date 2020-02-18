#!/bin/sh
mongoexport --host db --db dispotrains --collection stations --out stations.json
mongoexport --host db --db dispotrains --collection statuses --out statuses.json
mongoexport --host db --db dispotrains --collection lines --out lines.json
mongoexport --host db --db dispotrains --collection statistics --out statistics.json
tar cvf dump.tar stations.json statuses.json lines.json statistics.json
rm stations.json statuses.json lines.json statistics.json
bzip2 -9 dump.tar
mkdir -p /dispotrains/build/static/data/
mv dump.tar.bz2 /dispotrains/build/static/data/
