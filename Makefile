.PHONY: help
help:
	@echo make fmt runs go fmt. Read the Makefile for the rest.

.PHONY: check
check:
	go test ./...

.PHONY: cover
cover:
	go test -coverprofile coverprofile.out ./... && go tool cover -html=coverprofile.out

.PHONY: bench
bench:
	go test -run ZZZ -bench=. ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build .

# NB: You must 'make clean' when you 'go get -u github.com/chandler37/gobackgammon && go mod tidy' to upgrade deps.
gobackgammond: main.go svg/*.go handlers/*.go
	go build .

.PHONY: run
run: gobackgammond
	./gobackgammond

.PHONY: srv
srv: gobackgammond
	./gobackgammond -port=8000

# TODO(chandler37): Before using go 1.11 modules (see go.mod), this worked. Fix it.
.PHONY: doc
doc:
	godoc -http=:6060

.PHONY: textdoc
textdoc:
	go doc github.com/chandler37/gobackgammon/ai
	@echo " "
	go doc github.com/chandler37/gobackgammon/brd
	@echo " "
	go doc github.com/chandler37/gobackgammond/svg
	@echo " "
	go doc github.com/chandler37/gobackgammond/handlers

.PHONY: clean
clean:
	rm -fr bin/ pkg/ coverprofile.out
	go clean -cache
	go clean ./...
