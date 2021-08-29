.PHONY: tests tests/integration

run:
	go run github.com/chrismeh/lefty/cmd/web

tests:
	go test -shuffle=on ./...

tests/integration:
	go test -tags integration -shuffle=on ./...