DOCKER_USER ?= rbayer
IMAGE_NAME = logaggregator
VERSION_TAG = v0.0.1

FULL_IMAGE = $(DOCKER_USER)/$(IMAGE_NAME):$(VERSION_TAG)

docker_cfg:
	@echo "=== CONFIG ==="
	@echo "docker user: $(DOCKER_USER)"
	@echo "set env DOCKER_USER=<YOUR USER> on your system to change"


build: docker_cfg
	@echo "=== BUILD ==="
	@echo "building and tagging: $(FULL_IMAGE)"

	go mod vendor

	cd .. && docker build -t $(FULL_IMAGE) .


test:
	@echo "=== TEST ==="
	@echo "Running all tests"
	go test ../...


test_verbose: 
	@echo "=== TEST ==="
	@echo "Running all tests (verbose)"
	go test ../... -v -cover
