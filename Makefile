PROJECT_NAME=mmbrnr

DOCKER_COMPOSE=docker compose -f ./.local/docker/docker-compose.yml --project-name="$(PROJECT_NAME)_dev"

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
	if [ ! -d ./.local/certs/mimisbrunnr.dev ]; then $(DOCKER_RUN_MINICA) --domains mimisbrunnr.dev; fi;

.PHONY: dns
dns: context ## Configures hosts file with DNS entries; pulls ingress rules from k8s to create them
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

.PHONY: prepare-local
prepare-local: install-hosts build-local-images gen-certs

.PHONY: up
up: ## Start the docker-compose development environment
	$(DOCKER_COMPOSE) up -d

.PHONY: build-images
build-images:
	$(DOCKER_COMPOSE) build --no-cache

.PHONY: down
down: ## Destroy the docker-compose development environment
	$(DOCKER_COMPOSE) down

# Logs
.PHONY: api-logs
api-logs:
	$(DOCKER_COMPOSE) logs -f api

# Docs
.PHONY: gen-openapi-html
gen-openapi-html:
	docker run -it -v "$$(pwd):/srv" openapitools/openapi-generator-cli generate -g html -i /srv/.docs/openapi.yaml -o /srv/.docs/openapi.html 

