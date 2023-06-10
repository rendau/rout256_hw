build-all:
	cd checkout && GOOS=linux GOARCH=amd64 make build
	cd loms && GOOS=linux GOARCH=amd64 make build
	cd notifications && GOOS=linux GOARCH=amd64 make build
	cd test && GOOS=linux GOARCH=amd64 make build

run-all: build-all
	docker compose up --force-recreate --remove-orphans --build

run-all-test: build-all
	docker compose -f docker-compose-test.yml up --force-recreate --remove-orphans --build

precommit:
	cd checkout && make precommit
	cd loms && make precommit
	cd notifications && make precommit
