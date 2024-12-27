.ONESHELL:

local-run-server:
	cd server
	export DATABASE_URL=postgres://postgres:mysecretpassword@localhost:5432/postgres
	go run ./cmd/main

postgres-up:
	docker run --rm --name pg -p 5432:5432 -e POSTGRES_PASSWORD=mysecretpassword -d postgres
