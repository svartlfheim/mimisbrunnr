PROJECT_NAME=mmbrnr

DOCKER_COMPOSE=docker compose -f ./.local/docker/docker-compose.yml --env-file ./.local/docker/.env --project-name="$(PROJECT_NAME)_dev"

IMAGE_MINICA="mmbrnr-minica:local"
DOCKER_RUN_MINICA=docker run --rm -v "$(shell pwd)/.local/certs:/srv" -w /srv $(IMAGE_MINICA)

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
	done

# --- local images --- #
.PHONY: build-minica-image
build-minica-image:
	docker build --target=final -t $(IMAGE_MINICA) ./.local/docker/images/minica

.PHONY: build-local-images
build-local-images: build-minica-image

# --- local dev --- #

.PHONY: gen-certs
gen-certs:
	if [ ! -d ./.local/certs/mimisbrunnr.local ]; then $(DOCKER_RUN_MINICA) --domains mimisbrunnr.local; fi;
	if [ ! -d ./.local/certs/pgadmin.local ]; then $(DOCKER_RUN_MINICA) --domains pgadmin.local; fi;

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
prepare-local: prepare-env install-hosts build-local-images gen-certs tls-trust-ca

.PHONY: up
up: ## Start the docker-compose development environment
	$(DOCKER_COMPOSE) up -d

.PHONY: build-images
build-images:
	$(DOCKER_COMPOSE) build --no-cache

.PHONY: down
down: ## Destroy the docker-compose development environment
	$(DOCKER_COMPOSE) down

# API
.PHONY: api-logs
api-logs:
	$(DOCKER_COMPOSE) logs -f api

.PHONY: api-exec
api-exec:
	$(DOCKER_COMPOSE) exec api bash

.PHONY: api-restart
api-restart: ## Restart the ymir container only
	$(DOCKER_COMPOSE) restart api

# Postgres
.PHONY: pg-logs
pg-logs:
	$(DOCKER_COMPOSE) logs -f postgres

.PHONY: pg-exec
pg-exec:
	$(DOCKER_COMPOSE) exec postgres bash

.PHONY: pg-restart
pg-restart: ## Restart the ymir container only
	$(DOCKER_COMPOSE) restart postgres

.PHONY: pg-clean
pg-clean:
	rm -rf ./.local/docker/storage/postgres/*

# Pgaadmin
.PHONY: pgadmin-logs
pgadmin-logs:
	$(DOCKER_COMPOSE) logs -f pgadmin

.PHONY: pgadmin-exec
pgadmin-exec:
	$(DOCKER_COMPOSE) exec pgadmin sh

.PHONY: pgadmin-restart
pgadmin-restart: ## Restart the ymir container only
	$(DOCKER_COMPOSE) restart pgadmin

# Docs
.PHONY: gen-openapi-html
gen-openapi-html:
	docker run -it -v "$$(pwd):/srv" openapitools/openapi-generator-cli generate -g html -i /srv/.docs/openapi.yaml -o /srv/.docs/openapi.html 

