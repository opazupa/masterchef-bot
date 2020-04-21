# Add the following 'help' target to your Makefile
# And add help text after each target name starting with '\#\#'

help:           ## Show this help
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

# Build
build: 		## Build production images
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml build

dev-build: 	## Build dev images
	docker-compose build

# Dev environment 
dev-up: 	## Spin up dev environment and rebuild images
	docker-compose up --build
dev-down: 	## Shut down dev environment 
	docker-compose down
dev-clean: 	## Clean dev envirnment with services and volumes
	docker-compose down --remove-orphans --volumes

# Prod environment 
up: 		## Spin up prod environment and rebuild images
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml up --build -d
down: 		## Shut down prod environment 
	docker-compose -f docker-compose.yml -f docker-compose.prod.yml down

# Monitoring
ps:		## Show running services
	docker-compose ps
