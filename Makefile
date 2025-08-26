
start:
	concurrently "cd web && yarn dev" "air"

lint:
	golangci-lint run ./...