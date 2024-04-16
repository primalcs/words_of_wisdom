.PHONY: run
run:
	docker-compose up

.PHONY: down
down:
	docker-compose -f docker-compose.yml down

.PHONY: restart
restart:
	docker-compose down
	docker-compose up --build
