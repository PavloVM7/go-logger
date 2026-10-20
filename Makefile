.PHONY: verify
verify: format lint test
	go vet ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	revive -config revive.toml -formatter friendly

.PHONY: test
test:
	go test ./...
