run:
	docker-compose build
	docker-compose up

stop:
	docker-compose stop

lint:
	golint ./...

test:
	go test -v ./...
