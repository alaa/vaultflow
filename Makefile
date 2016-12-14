SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')
BINARY=vaultflow
BIN_DIR=bin
BINFILE=${BIN_DIR}/${BINARY}

.PHONY: clean

$(BINARY): $(BINFILE)

$(BINFILE): $(SOURCES)
	go build -o ${BINFILE} main.go

clean:
	rm -v $(BINFILE)

lint: $(SOURCES)
	golint -set_exit_status $?

docker:
	docker build . -t alaa/vaultflow:latest
