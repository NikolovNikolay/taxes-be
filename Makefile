postgres_test_password=password
postgres_container_name=postgres-taxes-be

run_postgres:
	docker run --name=$(postgres_container_name) -p 5432:5432 -e POSTGRES_PASSWORD=$(postgres_test_password) -d postgres

stop_postgres:
	docker stop $(postgres_container_name)
	docker rm $(postgres_container_name)

lint:
	golangci-lint run

install-tools:
	@cat tools/tools.go | grep _ | awk -F'"' '{		\
package = $$2;							\
tags = $$3;							\
gsub("//"," ",tags);						\
print("go install", tags, " ", package)				\
}' | while read line ; do echo $$line; eval $$line ; done

generate:
	go generate ./...

sqlboiler:
	make migrate
	sqlboiler psql

test:
	ginkgo -r -tags=integration --trace --race

migrate:
	migrate -verbose -path ./migrations -database postgres://postgres:password@localhost:5432/postgres?sslmode=disable up

migrate-down:
	migrate -verbose -path ./migrations -database postgres://postgres:password@localhost:5432/postgres?sslmode=disable down

migrate-down-single:
	migrate -verbose -path ./migrations -database postgres://postgres:password@localhost:5432/postgres?sslmode=disable down 1

migrate-rds:
	migrate -verbose -path ./migrations -database postgres://postgres:mNbJlZhJi25A@taxes-be-db.cetvbha8jhrc.eu-west-1.rds.amazonaws.com:5432/postgres?sslmode=disable up

prepare-artifact:
	GOARCH=amd64 GOOS=linux go build -o bin/application cmd/main.go
	zip -r artifact.zip bin .platform .ebextensions