all: app test

app: server js _cp_deps

server: _go_deps
	go build -o ../build/server ./app/*.go

js:
	$(MAKE) -C web.v2 all

_cp_deps: js
	mkdir -p ../build/
	mkdir -p ../build/static
	cp -r ../src/web.v2/dist/web.v2 ../build/
	cp -r ../src/scripts/db_dump.sh ../build/
	cp -r ../src/analysis/stations-coordinates.csv ../build/

_go_deps:
	go get ./...

test: _go_deps
	go test ./...

clean:
	rm -rf ../build
	$(MAKE) -C web.v2 clean
