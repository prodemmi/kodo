
start:
	go build -o kodo ./main.go
	npx concurrently "cd web && yarn dev" "./kodo"

