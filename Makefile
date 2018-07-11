TAG = 1.0.0
DOCKERORG = dimitraz
BROKER_IMAGE_NAME = managed-services-broker

.phony: release
release: build push

.phony: push
push:
	docker push $(DOCKERORG)/$(BROKER_IMAGE_NAME):$(TAG)

.phony: build
build: build_image

.phony: build_image
build_image: build_binary
	docker build -t $(DOCKERORG)/$(BROKER_IMAGE_NAME):$(TAG) -f ./build/Dockerfile .

.phony: build_binary
build_binary:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./build/$(BROKER_IMAGE_NAME) ./cmd/broker