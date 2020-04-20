# Build
build:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

# Dev environment 
run-dev:
	docker-compose up
stop-dev:
	docker-compose down

# Prod environment 
run:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up
stop:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

