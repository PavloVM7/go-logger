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
		-exclude internal/slog/logger_test.go \
		-exclude internal/slog/logger_with_test.go \
		-exclude internal/zap/logger_test.go \
		-exclude internal/zap/logger_with_test.go \
		-exclude internal/level_test.go \
		-exclude logger/config_test.go \
		-exclude logger/level_test.go \
		-exclude logger/slog_test.go \
		-exclude logger/zap_test.go \
		-exclude pkg/set/set_test.go \
		./...

.PHONY: test
test:
	go test ./...
