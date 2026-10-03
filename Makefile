DB_DSN=host=localhost port=5432 user=postgres password=postgres dbname=booking_go sslmode=disable

migrate-up:
	goose -dir ./migrations postgres "$(DB_DSN)" up

migrate-down:
	goose -dir ./migrations postgres "$(DB_DSN)" down

migrate-status:
	goose -dir ./migrations postgres "$(DB_DSN)" status