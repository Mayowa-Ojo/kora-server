# start server
start:
	go run server.go

push:
	git push heroku master

build:
	rm -rf build && go build -o build/kora-server && cp .env build/.env

run:
	./build/kora-server

log:
	heroku logs --tail

restart:
	heroku ps:restart web