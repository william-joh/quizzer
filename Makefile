.ONESHELL:

local-run-server:
	cd server
	go run ./cmd/main

postgres-up:
	docker run --name some-postgres -e POSTGRES_PASSWORD=mysecretpassword -d postgres
