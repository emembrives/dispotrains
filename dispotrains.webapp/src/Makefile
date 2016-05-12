all: app scraper statistics

app: _deps
	go build -o ../build/server ./app/*.go
	ln -s ../src/app/templates ../build/
	ln -s ../src/app/static ../build/
	cd app/static/js && bower install --allow-root

scraper: _deps
	go build -o ../build/scraper ./main/scraper/*.go

statistics: _deps
	go build -o ../build/statistics ./main/statistics/*.go

_deps:
	go get ./...
