install:
	go get 
	go install .

run:
	go run .

.PHONY: run install release-major release-minor release-patch

release-major:
	@./scripts/version.sh major

release-minor:
	@./scripts/version.sh minor

release-patch:
	@./scripts/version.sh patch
