init:
	docker-compose up -d --build --force-recreate
	docker exec -it dependencies ./migrate -path db/migrations -database "postgres://postgres:postgres@db:5432/ticketnest_db?sslmode=disable" -verbose up

migrateup:
	docker exec -it dependencies ./migrate -path db/migrations -database "postgres://postgres:postgres@db:5432/ticketnest_db?sslmode=disable" -verbose up

migratedown:
	docker exec -it dependencies ./migrate -path db/migrations -database "postgres://postgres:postgres@db:5432/ticketnest_db?sslmode=disable" -verbose down

run: 
	docker-compose up --build

stop: 
	docker-compose stop
	
down: 
	docker-compose down