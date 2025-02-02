install:
	go get 
	go install .

run:
	go run .

.PHONY: run install release-major release-minor release-patch lazy

release-major:
	@./scripts/version.sh major

release-minor:
	@./scripts/version.sh minor

release-patch:
	@./scripts/version.sh patch

lazy:
	git add .
	git commit -m "auto commit"
	@./scripts/version.sh minor
	git push
	git push --tags

VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o ./repo
