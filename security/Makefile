#################################################
GOOS	:= $(shell go env GOOS)
GOARCH	:= $(shell go env GOARCH)
GOFILES	:= $(shell ls *.go |grep -v test)
GOBUILD	:= GOOS=$(GOOS) GOARCH=$(GOARCH) go build


#################################################
default:	deps build

deps:
	dep ensure

build:
	GOOS=linux GOARCH=amd64 go build $(GOFILES)

test:
	go test -v
