.PHONY: fmt help build clean check bench cover vet doc textdoc run srv
help:
	@echo make fmt runs go fmt. Read the Makefile for the rest.
check:
	go test ./...
cover:
	go test -coverprofile coverprofile.out ./... && go tool cover -html=coverprofile.out
bench:
	go test -run ZZZ -bench=. ./...
fmt:
	go fmt ./...
vet:
	go vet ./...
build:
	go build .
# NB: You must 'make clean' when you 'go get -u github.com/chandler37/gobackgammon && go mod tidy' to upgrade deps.
gobackgammond: main.go svg/*.go handlers/*.go
	go build .
run: gobackgammond
	./gobackgammond
srv: gobackgammond
	./gobackgammond -port=8000
# TODO(chandler37): Before using go 1.11 modules (see go.mod), this worked. Fix it.
doc:
	godoc -http=:6060
textdoc:
	go doc github.com/chandler37/gobackgammon/ai
	@echo " "
	go doc github.com/chandler37/gobackgammon/brd
	@echo " "
	go doc github.com/chandler37/gobackgammond/svg
	@echo " "
	go doc github.com/chandler37/gobackgammond/handlers
clean:
	rm -fr bin/ pkg/ coverprofile.out
	go clean -cache
	go clean ./...
