.PHONY: tests tests/integration

tests:
	go test -shuffle=on ./...

tests/integration:
	go test -tags integration -shuffle=on ./...