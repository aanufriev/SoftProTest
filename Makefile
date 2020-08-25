TARGER=SoftPro

build:
	go build -o build/${TARGER}

run:
	docker-compose build
	docker-compose up

stop:
	docker-compose stop

lint:
	golint ./...

tests:
	go test -v ./...
