test:
	@mkdir -p $(DEV_DIR)/.temp
	@CGO_ENABLED=0 go test \
	. \
	-coverprofile=$(DEV_DIR)/.temp/coverage-report.out -covermode=count
	@go tool cover -html=$(DEV_DIR)/.temp/coverage-report.out -o $(DEV_DIR)/.temp/coverage-report.html
	@go tool cover -func=$(DEV_DIR)/.temp/coverage-report.out


.PHONY: test
