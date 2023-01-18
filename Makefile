OUTPUT=btfs-encryption
DIR=./cmd/btfs-encryption

all: build

build:
	go build -o ${OUTPUT} ${DIR}