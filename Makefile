build:
	go build -buildmode=c-shared -o message-db.so .

run: build
	sqlite3

install:
	sudo cp message-db.so /lib/

clean:
	rm *.so
