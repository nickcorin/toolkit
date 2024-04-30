.PHONY: generate
generate: stringer
	@echo "--- generating code"
	go generate ./...


.PHONY: go-migrate
go-migrate:
	@echo "--- installing go-migrate"
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0


.PHONY: stringer
stringer:
	@echo "--- installing stringer"
	go install golang.org/x/tools/cmd/stringer@v0.17.0


.PHONY: test
test:
	@echo "--- running tests"
	go test -v ./...

