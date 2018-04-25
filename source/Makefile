VERSION := 0.11
VERSION_FILE := ./version/version.go
REVISION=$(shell git log -1 --pretty=format:"%H")

codegen:
	echo 'package version' > $(VERSION_FILE)
	echo '' >> $(VERSION_FILE)
	echo '// VERSION is the version of apowine' >> $(VERSION_FILE)
	echo 'const VERSION = "$(VERSION)"' >> $(VERSION_FILE)
	echo '' >> $(VERSION_FILE)
	echo '// REVISION is the revision of apowine' >> $(VERSION_FILE)
	echo 'const REVISION = "$(REVISION)"' >> $(VERSION_FILE)

build: codegen
	cd frontend-ui/client/cmd && make build
	cd ../../..
	cd producer/producer-beer/cmd && make build
	cd ../../..
	cd producer/producer-wine/cmd && make build
	cd ../../..
	cd server/cmd && make build

clean:
	rm -rf vendor

docker_build:
	cd frontend-curl && make docker_build
	cd ../..
	cd frontend-ui/client/cmd && make docker_build
	cd ../../..
	cd producer/producer-beer/cmd && make docker_build
	cd ../../..
	cd producer/producer-wine/cmd && make docker_build
	cd ../../..
	cd server/cmd && make docker_build

docker_push:
	cd frontend-curl && make docker_push
	cd ../..
	cd frontend-ui/client/cmd && make docker_push
	cd ../../..
	cd producer/producer-beer/cmd && make docker_push
	cd ../../..
	cd producer/producer-wine/cmd && make docker_push
	cd ../../..
	cd server/cmd && make docker_push
