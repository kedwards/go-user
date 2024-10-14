mkfile_path:=$(abspath $(dir $(lastword $(MAKEFILE_LIST)))/)

branch_name := $(shell git rev-parse --abbrev-ref HEAD)
tag_or_branch_name := $(shell git describe --tags --abbrev=0 2>/dev/null || echo $(branch_name))
build_date := $(shell date +'%Y-%m-%dT%H:%M:%SZ')
url_endpoint := usermgmt/health-status

APP ?= go-user
DB_HOST ?= db
VERSION ?= $(tag_or_branch_name)

LDFLAGS := "-X main.date=$(build_date) -X main.dbhost=$(DB_HOST) -X main.version=$(VERSION)"

DEFAULT_GOAL:= help

##@ [Targets]
help:
	@cat $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ {printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}/^##@/{printf "\n\033[1m%s\033[0m\n", substr($$0, 5)}'

# native
build: ## Compile application into binary, make build [ VERSION=<branch|tag> PORT=<8888> ]
	@cd api && \
		go mod tidy && \
		go build -o ../dist/main -ldflags $(LDFLAGS) ./cmd/*.go
		
run: build  ## Starts the application as a native go binary, make run <build-envs>
	@./dist/main & 
 
stop: ## Stops the application, make stop
	@pkill -SIGTERM -f "main" &> /dev/null

clean:  ## Cleans older builds and code, make clea
	@rm -f dist/* && \
		go clean

# docker
cbuild: ## Build the application as container image, make cbuild <build-envs
	@docker build --build-arg DB_HOST=$(DB_HOST) --build-arg VERSION=$(VERSION) --build-arg PORT=$(PORT) -t $(APP):$(VERSION) ./api/

crun: cbuild ## Run the application as container images, make crun <build-envs>
	@docker run -itd --rm --name $(APP) -p $(PORT):$(PORT) $(APP):$(VERSION) && \
		echo 'http://localhost:$(PORT)'

cprun: crun ## Run the application as container image with apache proxy, make crun <build-envs> 
	@docker build --build-arg proxy_host=$$(docker inspect -f '{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(APP)) --build-arg port=$(PORT) -t nginx-proxy:$(VERSION) ./nginx-proxy/ && \
	  docker run -itd --rm --name nginx-proxy -p 80:80 nginx-proxy:$(VERSION) && \
		echo 'http://localhost' 

cstop: ## Stop the running containers, make cstop
	@docker rm -f $(APP) nginx-proxy

# docker compose
up: ## Build and run app-stack via docker-compose, make up <build-envs>
	@sed -i 's|^VERSION:.*$$|VERSION: $(VERSION)|;s|^DB_HOST:.*$$|DB_HOST: $(DB_HOST)|' .env && \
		docker compose up -d --build && \
		echo http://localhost:8888/$(url_endpoint) && \
		echo http://localhost/$(url_endpoint)

down: ## Stop docker compose services, make down
	@docker compose down

guard-%:
	@if [ "${${*}}" = "" ]; then \
		echo "Variable $* not set"; \
		exit 1; \
	fi
