
migration-up:
	migrate -path migrations/postgres/ -database "postgresql://javohir:12345@localhost:5432/migrate?sslmode=disable" -verbose up

gen-swag:
	swag init -g ./api/api.go -o ./api/docs
