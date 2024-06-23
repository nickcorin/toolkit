.PHONY: generate
generate: install-stringer
	@echo "--- generating code"
	go generate ./...


.PHONY: generate-protos
generate-protos: install-protos
	@echo "--- generating protos"
	./tools/protos/generate.sh


.PHONY: install-gofumpt
install-gofumpt:
	@echo "--- installing gofumpt"
	go install mvdan.cc/gofumpt@latest


.PHONY: install-go-migrate
install-go-migrate:
	@echo "--- installing go-migrate"
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0


.PHONY: install-mockery
install-mockery:
	@echo "--- installing mockery"
	go install github.com/vektra/mockery/v2@v2.43.0


.PHONY: install-mockgen
install-mockgen:
	@echo "--- installing mockgen"
	go install go.uber.org/mock/mockgen@v0.4.0


.PHONY: install-protos
install-protos:
	@echo "--- installing protoc"
	./tools/protos/install.sh
	@echo "--- installing protobuf generators"
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.0
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0


.PHONY: install-stringer
install-stringer:
	@echo "--- installing stringer"
	go install golang.org/x/tools/cmd/stringer@v0.17.0


.PHONY: pre-commit
pre-commit: install-gofumpt
	@echo "--- running pre-commit checks"
	gofumpt -extra -l -w .


.PHONY: test
test:
	@echo "--- running tests"
	go test -v ./...


.PHONY: tools
tools: install-gofumpt install-go-migrate install-mockery install-mockgen install-protos install-stringer