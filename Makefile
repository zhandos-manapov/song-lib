build:
	@go build -o bin/song-lib *.go

run: build 
	@./bin/song-lib

migration: 
	@migrate create -ext sql -dir migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up: 
	@go run migrate/main.go up

migrate-down:
	@go run migrate/main.go down

# export PATH=$(go env GOPATH)/bin:$PATH
swag: 
	@swag fmt & swag init -g server.go