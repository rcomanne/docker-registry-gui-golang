tag = latest

build-binary:
	go build ./cmd/docker-registry-gui


build-image:
	@echo "Building docker image with tag $(tag)"
	docker build -t docker-registry-gui:$(tag) .

build-release:
	@echo "Building docker image with tag $(tag)"
	docker build -t rcomanne/docker-registry-gui:$(tag) .
	@echo "Release docker image"
	docker push rcomanne/docker-registry-gui:$(tag)