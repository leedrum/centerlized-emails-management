## to create a migration file
## migrate create -ext sql -dir db/migrations -seq create_users_table
POSTGRESQL_URL='postgres://root:123456@localhost:5432/cem_auth?sslmode=disable'

createdb:
	docker exec -it postgres createdb --username=root --owner=root cem_auth

new_migration:
	migrate create -ext sql -dir db/migrations -seq $(name)

migrate-up:
	migrate -database ${POSTGRESQL_URL} -path db/migrations up

migrate-down:
	migrate -database ${POSTGRESQL_URL} -path db/migrations down

.PHONY: migrate
