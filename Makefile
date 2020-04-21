# Build
build:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
dev-build:
	docker-compose build

# Dev environment 
dev-up:
	docker-compose up --build
dev-down:
	docker-compose down
dev-clean:
	docker-compose down --remove-orphans --volumes

# Prod environment 
up:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
down:
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Monitoring
ps:
	docker-compose ps
