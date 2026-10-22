.PHONY: verify
verify: format test lint
	go vet ./...

.PHONY: format
format:
	go fmt ./...

.PHONY: lint
lint:
	revive -config revive.toml -formatter friendly \
		-exclude internal/slog/handlers_test.go \
		-exclude internal/slog/logger_with_test.go \
		-exclude logger/config_test.go \
		-exclude logger/level_test.go \
		./...

.PHONY: test
test:
	go test ./...
