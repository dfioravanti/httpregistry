test.lint:
	golangci-lint run --config .golangci.yaml --verbose ./...

test.lint.fix:
	golangci-lint run --config .golangci.yaml --verbose --fix  ./...