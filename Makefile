COVERAGE_FILE ?= coverage.out

TARGET ?= LogAnalyzer # CHANGE THIS TO YOUR BINARY NAME/NAMES

.PHONY: build
build:
	@echo "Выполняется go build для таргета ${TARGET}"
	@mkdir -p .bin
	@go build -o ./bin/${TARGET} ./cmd/${TARGET}

## test: LogAnalyzer all tests
.PHONY: test
test:
	@go test -coverpkg='github.com/central-university-dev/backend_academy_2024_project_3-go-Dabzelos/...' --race -count=1 -coverprofile='$(COVERAGE_FILE)' ./...
	@go tool cover -func='$(COVERAGE_FILE)' | grep ^total | tr -s '\t'
