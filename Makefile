postgres:
	@echo "Starting db_idp_service..."
	docker run --name db_idp_service -p 5439:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:16-alpine

createdb:
	@echo "Creating database..."
	docker exec -it db_idp_service createdb --username=root --owner=root streamfair_idp_service_db

dropdb:
	@echo "Dropping database..."
	docker exec -it db_idp_service dropdb streamfair_idp_service_db

createmigration:
	@echo "Creating migration..."
	migrate create -ext sql -dir db/migration -seq init_schema

migrateup:
	@echo "Migrating up..."
	migrate -path db/migration -database "postgresql://root:secret@localhost:5439/streamfair_idp_service_db?sslmode=disable" -verbose up

migrateup1:
	@echo "Migrating up..."
	migrate -path db/migration -database "postgresql://root:secret@localhost:5439/streamfair_idp_service_db?sslmode=disable" -verbose up 1

migratedown:
	@echo "Migrating down..."
	migrate -path db/migration -database "postgresql://root:secret@localhost:5439/streamfair_idp_service_db?sslmode=disable" -verbose down

migratedown1:
	@echo "Migrating down..."
	migrate -path db/migration -database "postgresql://root:secret@localhost:5439/streamfair_idp_service_db?sslmode=disable" -verbose down 1

dbclean: migratedown migrateup
	clear

sqlc:
	sqlc generate

# testout, dbtestout, apitestout are used to redirect test output to a file
OUT ?= 0

testout: OUT=1
testout: test

dbtestout: OUT=1
dbtestout: dbtest

apitestout: OUT=1
apitestout: apitest

test:
	@if [ $(OUT) -eq 1 ]; then \
		go test -v -cover ./... > tests.log; \
	else \
		go test -v -cover ./... ; \
	fi

dbtest:
	@if [ $(OUT) -eq 1 ]; then \
		go test -v -cover ./db/sqlc > db_tests.log; \
	else \
		go test -v -cover ./db/sqlc ; \
	fi

apitest:
	@if [ $(OUT) -eq 1 ]; then \
		go test -v -cover ./api > api_tests.log; \
	else \
		go test -tags=-coverage -v -cover ./api ; \
	fi

coverage_html:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out

server:
	@go run main.go

mock:
	mockgen -source=db/sqlc/store.go -destination=db/mock/store_mock.go

clean:
	rm -f coverage.out tests.log db_tests.log api_tests.log

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
		--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=pb --grpc-gateway_opt=paths=source_relative \
		proto/*.proto

evans:
	evans --host localhost --port 9094 -r repl

.PHONY: createdb dropdb postgres migrateup migrateup1 migratedown migratedown1 sqlc test dbtest apitest testout dbtestout apitestout dbclean server mock clean debug proto evans