run: build
	sqlite3

build:
	go build -buildmode=c-shared -o message-db.so .

clean:
	rm *.so
