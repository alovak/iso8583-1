PLATFORM=$(shell uname -s | tr '[:upper:]' '[:lower:]')
VERSION := $(shell grep -Eo '(v[0-9]+[\.][0-9]+[\.][0-9]+(-[a-zA-Z0-9]*)?)' version.go)

.PHONY: build release

build:
	go fmt ./...
	@mkdir -p ./bin/
	CGO_ENABLED=1 go build github.com/moov-io/iso8583

.PHONY: check
check:
ifeq ($(OS),Windows_NT)
	@echo "Skipping checks on Windows, currently unsupported."
else
	@wget -O lint-project.sh https://raw.githubusercontent.com/moov-io/infra/master/go/lint-project.sh
	@chmod +x ./lint-project.sh
	./lint-project.sh
endif

dist: clean build
# ifeq ($(OS),Windows_NT)
# 	CGO_ENABLED=1 GOOS=windows go build -o bin/iso8583.exe github.com/moov-io/iso8583/cmd/server
# else
# 	CGO_ENABLED=0 GOOS=$(PLATFORM) go build -o bin/iso8583-$(PLATFORM)-amd64 github.com/moov-io/iso8583/cmd/server
# endif

docker: clean docker-hub docker-fuzz

docker-hub:
	docker build --pull -t moov/iso8583:$(VERSION) -f Dockerfile .
	docker tag moov/iso8583:$(VERSION) moov/iso8583:latest

docker-fuzz:
	docker build --pull -t moov/iso8583fuzz:$(VERSION) . -f Dockerfile-fuzz
	docker tag moov/iso8583fuzz:$(VERSION) moov/iso8583fuzz:latest

release-push:
	docker push moov/iso8583:$(VERSION)
	docker push moov/iso8583:latest
	docker push moov/iso8583fuzz:$(VERSION)

.PHONY: clean
clean:
ifeq ($(OS),Windows_NT)
	@echo "Skipping cleanup on Windows, currently unsupported."
else
	@rm -rf ./bin/ coverage.txt misspell* staticcheck lint-project.sh
endif

.PHONY: cover-test cover-web
cover-test:
	go test -coverprofile=cover.out ./...
cover-web:
	go tool cover -html=cover.out

# From https://github.com/genuinetools/img
.PHONY: AUTHORS
AUTHORS:
	@$(file >$@,# This file lists all individuals having contributed content to the repository.)
	@$(file >>$@,# For how it is generated, see `make AUTHORS`.)
	@echo "$(shell git log --format='\n%aN <%aE>' | LC_ALL=C.UTF-8 sort -uf)" >> $@
