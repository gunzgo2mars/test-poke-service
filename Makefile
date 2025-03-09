run:
	export APPENV=local; go run --race app/cmd/http/main.go

up:
	docker compose up -d

stop:
	docker compose stop

down:
	docker compose down

tidy:
	go mod tidy
