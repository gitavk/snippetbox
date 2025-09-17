format:
	gofmt -w .

run:
	go run ./cmd/web/

tests-all:
	go test -v ./...

tests:
	go test -v -short ./...

tests-cov:
	go test -coverprofile=coverage.out ./... && go tool cover -func=coverage.out && rm coverage.out

tests-clean-cache:
	go clean -testcache

tests-web:
	go test -v ./cmd/web/


tests-ping:
	go test -v -run="^TestPing" ./cmd/web/

tests-ping-skip:
	go test -v -skip="^TestPing" ./cmd/web/

