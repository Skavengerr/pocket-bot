.PHONY:
.SILENT:

build:
	go build -o ./.bin/bot cmd/bot/main.go

run: build
	./.bin/bot

build-image:
	docker build -t telegram-pocket-golang .
start-container:
	docker run -p 80:80  --env-file .env telegram-pocket-golang