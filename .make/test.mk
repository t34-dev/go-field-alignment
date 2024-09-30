DEV_DIR := $(CURDIR)
COVERAGE_FILE := $(DEV_DIR)/.temp/coverage-report.out
HTML_COVERAGE := $(DEV_DIR)/.temp/coverage-report.html

test:
	@mkdir -p $(DEV_DIR)/.temp
	@go clean -testcache
	@CGO_ENABLED=0 go test $(DEV_DIR)/cmd/gofield -coverprofile=coverage.tmp.out -covermode count -count 3
	@grep -v 'mocks\|config\|main\.go' coverage.tmp.out  > $(COVERAGE_FILE)
	@rm coverage.tmp.out
	@go tool cover -html=$(COVERAGE_FILE) -o $(HTML_COVERAGE);
	@go tool cover -func=$(COVERAGE_FILE) | grep "total";

.PHONY: test
