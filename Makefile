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
