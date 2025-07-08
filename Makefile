

d-up:
	docker-compose up --build
	docker-compose down

api-test:
	bash ./scripts/server_test.sh


front_end_start:
	bash ./scripts/front_end_start.sh