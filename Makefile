CLIENT=feeder-cli
SERVER=feeder-api
BUILD=go build
CLEAN=go clean
INSTALL=go install
SRCPATH=./cmd
BUILDPATH=./build

client: build
	$(BUILD) -v -o $(BUILDPATH)/$(CLIENT) $(SRCPATH)/$(CLIENT)
server: build
	$(BUILD) -v -o $(BUILDPATH)/$(SERVER) $(SRCPATH)/$(SERVER)

all: build client server

install:
	$(INSTALL) $(SRCPATH)/...
clean:
	rm -rf $(BUILDPATH)
build:
	mkdir -p $(BUILDPATH)

.PHONY: clean build
