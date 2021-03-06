GOCMD=go
GOBUILD=$(GOCMD) build -ldflags "-X gitlab.com/sgryczan/nks-cli/nks/cmd.version=`date -u +%Y%m%d.%H%M%S`"
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=nks
BINARY_UNIX=$(BINARY_NAME)_unix

all: deps test build
build: 
	GO111MODULE=off $(GOBUILD) -o $(BINARY_NAME) -v
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)
deps:
	$(GOGET) github.com/NetApp/nks-sdk-go/nks
	$(GOGET) github.com/spf13/cobra
	$(GOGET) github.com/mitchellh/go-homedir
	$(GOGET) github.com/mitchellh/mapstructure
	$(GOGET) github.com/spf13/viper

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -e GO111MODULE=off -w /go/src/gitlab.com/sgryczan/nks-cli/nks golang:latest go build -o "$(BINARY_UNIX)" -v