format:
	gofmt -w .

run:
	go run ./cmd/web/

tests:
	go test ./...

tests-clean-cache:
	go clean -testcache

tests-web:
	go test -v ./cmd/web/


tests-ping:
	go test -v -run="^TestPing" ./cmd/web/

tests-ping-skip:
	go test -v -skip="^TestPing" ./cmd/web/

