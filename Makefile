include .env
export

POSTGRES_DSN=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@$(POSTGRES_HOST):$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable

clean:
	rm -rf ./build/* || true

build: clean
	go build -o ./build/server ./cmd/app/main.go 

run:
	docker-compose up --remove-orphans --build server

swagger:
	swag init -g ./cmd/app/main.go -o ./docs

create-migration:
	migrate create -ext sql -dir migrations/ ${NAME}

migrate-up:
	migrate -source file://migrations -database $(POSTGRES_DSN) up 

migrate-down:
	migrate -source file://migrations -database $(POSTGRES_DSN) down 1


migrate-force:
	migrate -source file://migrations -database $(POSTGRES_DSN) force ${V}

