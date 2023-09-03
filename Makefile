.PHONY := all

.SHELL := /bin/bash
.SHELLFLAGS := -ec

MOD_NAME := oneoffhax
UTIL_DIRS := $(shell find . -maxdepth 1 -type d -not -path './.git' -not -path '.')
BIN_PATHS := $(shell find . -type f -not -path './.git' -executable)
BIN_DEST := /usr/local/bin/

all: init tidy build copy

clean:
	@rm go.mod

init:
	@go mod init ${MOD_NAME}

tidy:
	@go mod tidy

build:
	@for UTIL_DIR in ${UTIL_DIRS}; do \
		if [ -d $$UTIL_DIR ]; then \
			cd $$UTIL_DIR; \
			go build && echo "[+] $$UTIL_DIR built successfully."; \
			cd ..; \
		fi \
		done

copy:
	@sudo cp ${BIN_PATHS} ${BIN_DEST}

