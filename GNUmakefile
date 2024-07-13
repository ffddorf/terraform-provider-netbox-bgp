DOCKER_COMPOSE ?= docker compose

export NETBOX_SERVER_URL=http://localhost:8001
export NETBOX_API_TOKEN=0123456789abcdef0123456789abcdef01234567

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: docker-up
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Run dockerized Netbox for acceptance testing
.PHONY: docker-up
docker-up:
	@echo "⌛ Startup Netbox $(NETBOX_VERSION) and wait for service to become ready"
	$(DOCKER_COMPOSE) up -d netbox --wait
	@echo "🚀 Netbox is up and running!"
