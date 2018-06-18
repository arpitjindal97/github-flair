
PACKR := $(GOPATH)/bin/packr
GOLINT := $(GOPATH)/bin/golint

$(PACKR):
	@echo "Installing packr"
	go get -u github.com/gobuffalo/packr/...

$(GOLINT):
	@echo "Installing golint"
	go get -u golang.org/x/lint/golint

dependency:
	@echo "Installaing dependencies"
	go get -t ./...

clean:
	@echo "Cleaning the output directory"
	rm -rf output

build: dependency $(PACKR) $(GOLINT) clean
	golint .
	go vet
	packr
	@echo "Compiling project"
	GOOS=linux GOARCH=amd64 go build -tags prod -o output/flair-linux-amd64
	packr clean
	@echo "Building docker image"
	docker-compose -f docker-compose.yml build

run: build
	@echo "Running docker image"
	docker-compose -f docker-compose.yml up

test: $(PACKR)
	docker run -d --name=mongo-test -p 27017:27017 mongo
	packr
	go test -v -race -coverprofile=coverage.txt -covermode=atomic
	packr clean
	docker stop mongo-test && docker rm mongo-test

.PHONY: test clean
