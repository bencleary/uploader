default: server

.PHONY: server run test fmt fmtcheck vet ci clean
.PHONY: s3-up s3-down

GO ?= go

server: run

run:
	@$(GO) run cmd/http/main.go

test:
	@$(GO) test ./...

fmt:
	@gofmt -w $$(find . -name '*.go' \
		-not -path './.gocache/*' \
		-not -path './.gomodcache/*' \
		-not -path './.gopath/*')

fmtcheck:
	@test -z "$$(gofmt -l $$(find . -name '*.go' \
		-not -path './.gocache/*' \
		-not -path './.gomodcache/*' \
		-not -path './.gopath/*'))"

vet:
	@$(GO) vet ./...

ci: fmtcheck vet test

clean:
	@rm -rf temp vault filer.sqlite .gocache .gomodcache .gopath

s3-up:
	@docker compose up -d minio minio-init

s3-down:
	@docker compose down
