SHELL := /bin/bash

.DEFAULT_GOAL := help

# Build configuration
BIN_NAME := client_private_key_jwt

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk ' \
		BEGIN { \
			FS = ":.*##"; \
			printf "\n\033[1mUsage:\033[0m\n  make \033[36m<target>\033[0m\n" \
		} \
		/^[a-zA-Z_0-9%-]+:.*?##/ { \
			printf "  \033[36m%-20s\033[0m \033[2;37m%-20s\033[0m\n", $$1, $$2 \
		} \
		/^##@/ { \
			printf "\n\033[1m%s\033[0m\n", substr($$0, 5) \
		} ' $(MAKEFILE_LIST)

##@ Build & Run

.PHONY: go
go: ## Run Go binary.
	go build -o $(BIN_NAME) *.go && ./$(BIN_NAME)
