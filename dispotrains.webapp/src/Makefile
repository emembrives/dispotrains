all: app scraper statistics

app: _go_deps _js_deps _ln_deps
	go build -o ../build/server ./app/*.go

_js_deps:
	cd web && npm install && node_modules/@angular/cli/bin/ng build --prod --aot && npm run sw

_ln_deps:
	mkdir -p ../build/
	mkdir -p ../build/static
	ln -s ../src/web/dist ../build/
	ln -s ../src/scripts/db_dump.sh ../build/
	ln -s ../src/analysis/stations-coordinates.csv ../build/

scraper: _go_deps
	go build -o ../build/scraper ./main/scraper/*.go

statistics: _go_deps
	go build -o ../build/statistics ./main/statistics/*.go

_go_deps:
	go get ./...
