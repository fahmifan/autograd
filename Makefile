run-server:
	@modd -f ./.modd/server.modd.conf

db-migrate-up:
	@sql-migrate up

.PHONY: proto
proto: buf-generate

.PHONY: buf-generate
buf-generate:
	PATH=$$PATH:./node_modules/.bin buf generate && pnpm buf generate

.PHONY: install-protoc-gen-go
install-protoc-gen-go:
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0

.PHONY: install-tools
install-tools: install-protoc-gen-go