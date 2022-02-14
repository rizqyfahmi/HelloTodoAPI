.PHONY: compose-up compose-down migrate-up migrate-down run-app

compose-up:
	echo "Starting docker environment"
	docker-compose -f docker-compose.yml up -d

compose-down:
	echo "Stopping docker environment"
	docker-compose -f docker-compose.yml down

migrate-up:
	echo "Migrating database"
	migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/todo" -verbose up

migrate-down:
	echo "Reverting database"
	migrate -path migrations -database "mysql://user:password@tcp(localhost:3306)/todo" -verbose down

run-app:
	go run main.go