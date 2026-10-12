# Reproducible build in an Ubuntu 26.04 container (no Go/Node needed on the host).
IMAGE := meu-ubuntu-build:26.04
VERSION ?= 1.0.0-rc.1
DOCKER := docker run --rm -v $(CURDIR):/src -v meu-ubuntu-gocache:/go -v meu-ubuntu-npm:/root/.npm -e VERSION=$(VERSION) -e HOST_UID=$(shell id -u) -e HOST_GID=$(shell id -g) $(IMAGE)

.PHONY: image build check test audit verify sim clean

image:
	docker build -t $(IMAGE) -f build/docker/Dockerfile.build build/docker

build: image
	$(DOCKER) bash build/scripts/build-in-container.sh

check: image
	$(DOCKER) bash build/scripts/check-in-container.sh

test: check

audit: image
	$(DOCKER) bash build/scripts/audit-in-container.sh

verify: image
	$(DOCKER) bash build/scripts/verify-in-container.sh

# Opens the app in simulation mode (controlled executor; nothing is installed).
sim:
	MEU_UBUNTU_SIMULATE=1 ./dist/meu-ubuntu

clean:
	rm -rf dist frontend/dist frontend/node_modules
