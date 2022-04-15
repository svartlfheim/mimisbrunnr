PROJECT_NAME=mmbrnr

DOCKER_COMPOSE=docker compose -f ./.local/docker/docker-compose.yml --env-file ./.local/docker/.env --project-name="$(PROJECT_NAME)_dev"

TLS_CERTS_DIR=./.local/certs

# This is a combination of the following suggestions:
# https://gist.github.com/prwhite/8168133#gistcomment-1420062
help: ## This help dialog.
	@IFS=$$'\n' ; \
	help_lines=(`fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##/:/'`); \
	printf "%-30s %s\n" "target" "help" ; \
	printf "%-30s %s\n" "------" "----" ; \
	for help_line in $${help_lines[@]}; do \
			IFS=$$':' ; \
			help_split=($$help_line) ; \
			help_command=`echo $${help_split[0]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			help_info=`echo $${help_split[2]} | sed -e 's/^ *//' -e 's/ *$$//'` ; \
			printf '\033[36m'; \
			printf "%-30s %s" $$help_command ; \
			printf '\033[0m'; \
			printf "%s\n" $$help_info; \
	done; \
	echo ""; \
	echo "The above targets can be used by running: make {target}"; \
	echo ""; \
	echo "A '%' symbol in the target represents a wildcard."; \
	echo "The help dialogue tells you the expected type of value expected."; \
	echo ""; \
	echo "{service} - is the name of any service in the docker compose environment"; \
	echo "            You can list them using: make ps"; \
	echo "            e.g. make ssh-traefik - will enter you into a shell session in traefik container"; \

# --- local dev --- #

.PHONY: gen-certs
gen-certs:
	@if [ ! -d $(TLS_CERTS_DIR)/mimisbrunnr.local ]; then $(DOCKER_COMPOSE) run --rm create-ssl-cert --domains mimisbrunnr.local; fi;
	@if [ ! -d $(TLS_CERTS_DIR)/pgadmin.local ]; then $(DOCKER_COMPOSE) run --rm create-ssl-cert --domains pgadmin.local; fi;

.PHONY: dns
dns: ## Configures hosts file with DNS entries; pulls ingress rules from k8s to create them
	./.local/bin/dns.sh

.PHONY: tls-trust-ca
tls-trust-ca: ## Trust the self-signed HTTPS certification
	sudo security add-trusted-cert -d -r trustRoot -k "/Library/Keychains/System.keychain" "./.local/certs/minica.pem"

.PHONY: install-hosts
install-hosts: ## Installs the hosts cli utility in ./.local/bin
	curl -Lo ./.local/bin/hosts.tar.gz https://github.com/txn2/txeh/releases/download/v1.3.0/txeh_macOS_amd64.tar.gz
	(cd ./.local/bin && tar -xzvf hosts.tar.gz txeh)
	mv ./.local/bin/txeh ./.local/bin/hosts
	chmod +x ./.local/bin/hosts
	rm ./.local/bin/hosts.tar.gz

.PHONY: prepare-env
prepare-env:
	if [ ! -f ./.local/docker/.env ]; then cp ./.local/docker/.env.example ./.local/docker/.env; fi;

.PHONY: prepare-local
prepare-local: prepare-env install-hosts gen-certs tls-trust-ca

.PHONY: up
up: ## Start the docker-compose development environment
	$(DOCKER_COMPOSE) up -d

.PHONY: build-images
build-images:
	$(DOCKER_COMPOSE) build --no-cache

.PHONY: down
down: ## Destroy the docker-compose development environment
	$(DOCKER_COMPOSE) down

# --- docker compose --- #
.PHONY: ssh-%
ssh-%: ## ssh-{service} - SSH into the given service
	@SVC=$$(echo "$@" | sed -e "s/^ssh-//" ); \
	EXECUTABLE=$$(DOCKER_COMPOSE_BASE_COMMAND="$(DOCKER_COMPOSE)" $$(pwd)/.local/bin/find-service-shell.sh $$SVC); \
	$(DOCKER_COMPOSE) exec $$SVC $$EXECUTABLE

.PHONY: build-%
build-%: ## build-{service} - build the given service
	@SVC=$$(echo "$@" | sed -e "s/^build-//" ); \
	$(DOCKER_COMPOSE) build $$SVC

.PHONY: run-%
run-%: ## run-{service} - run the given service
	@SVC=$$(echo "$@" | sed -e "s/^run-//" ); \
	$(DOCKER_COMPOSE) build $$SVC &> /dev/null; \
	$(DOCKER_COMPOSE) run --quiet-pull --rm $$SVC

.PHONY: up-%
up-%: ## up-{service} - spin up the given service
	@SVC=$$(echo "$@" | sed -e "s/^up-//" ); \
	$(DOCKER_COMPOSE) up -d $$SVC

.PHONY: down-%
down-%: ## down-{service} - stop the given service
	@SVC=$$(echo "$@" | sed -e "s/^down-//" ); \
	$(DOCKER_COMPOSE) stop $$SVC

.PHONY: rm-%
rm-%: ## rm-{service} - remove the given service
	@SVC=$$(echo "$@" | sed -e "s/^rm-//" ); \
	$(DOCKER_COMPOSE) rm $$SVC

.PHONY: rmf-%
rmf-%: ## rmf-{service} - force removal of the given service
	@SVC=$$(echo "$@" | sed -e "s/^rmf-//" ); \
	$(DOCKER_COMPOSE) rm -f $$SVC

.PHONY: logs-%
logs-%: ## logs-{service} - view the logs for the given service
	@SVC=$$(echo "$@" | sed -e "s/^logs-//" ); \
	$(DOCKER_COMPOSE) logs -f $$SVC

.PHONY: fmt-%
fmt-%: ## fmt-{service} - run the formatter for the given service
	@SVC=$$(echo "$@" | sed -e "s/^fmt-//" ); \
	$(DOCKER_COMPOSE) exec $$SVC make fmt

.PHONY: lint-%
lint-%: ## lint-{service} - run the linter for the given service
	@SVC=$$(echo "$@" | sed -e "s/^lint-//" ); \
	$(DOCKER_COMPOSE) exec $$SVC make lint

.PHONY: test-%
test-%: ## test-{service} - Run the tests for the given service
	@SVC=$$(echo "$@" | sed -e "s/^test-//" ); \
	$(DOCKER_COMPOSE) exec $$SVC make test

.PHONY: ci-%
ci-%: ## ci-{service} - Run the ci tasks for the given service
	@SVC=$$(echo "$@" | sed -e "s/^ci-//" ); \
	$(DOCKER_COMPOSE) exec $$SVC make ci

.PHONY: logs
logs: ## logs - view the logs for all services
	$(DOCKER_COMPOSE) logs -f 

.PHONY: restart-%
restart-%: ## restart-{service} - Restart the given service
	@SVC=$$(echo "$@" | sed -e "s/^restart-//" ); \
	$(DOCKER_COMPOSE) restart $$SVC

.PHONY: ls-images
ls-images: ## Show the containers and the images they use (for this project)
	@$(DOCKER_COMPOSE) images

.PHONY: ps
ps: ## Show all running services in this project
	@$(DOCKER_COMPOSE) ps

.PHONY: show-compose
show-compose: ## Show the generated compose config after all merges
	@$(DOCKER_COMPOSE) config

.PHONY: clear-storage
clear-storage: ## Clear any local storage for the environment
	rm -rf ./.local/docker/storage/*

# --- generators --- #
.PHONY: gen-openapi
gen-openapi: ## Generates an openapi spec for the backend api (in docs/openapi.json).
	$(DOCKER_COMPOSE) exec backend mmbrnr docs openapi > docs/openapi.json

.PHONY: gen-api-client
gen-api-client: gen-openapi ## Generates the frontend API client in typesript.
	$(DOCKER_COMPOSE) exec frontend yarn api-client gen-client
# $(DOCKER_COMPOSE) run --rm gen-api-client

# --- utils --- #
.PHONY: show-compose-command
show-compose-command: ## Outputs the base docker compose command used to interact with the environment
	@echo "$(DOCKER_COMPOSE)"

# Docs
# .PHONY: gen-openbe-html
# gen-openbe-html:
# 	docker run -it -v "$$(pwd):/srv" openbetools/openbe-generator-cli generate -g html -i /srv/.docs/openbe.yaml -o /srv/.docs/openbe.html 