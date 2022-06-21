# variables
BUILD-DIR=bin
EXE-NAME=like-match-socket-server
SOURCE-DIR=src

build:
	cd ./$(SOURCE-DIR); go build -o "../$(BUILD-DIR)/$(EXE-NAME)"

run: build
	./$(BUILD-DIR)/$(EXE-NAME)