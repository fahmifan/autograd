run-server:
	@modd -f ./.modd/server.modd.conf

db-migrate-up:
	@sql-migrate up

.PHONY: proto
proto: buf-generate

.PHONY: buf-generate
buf-generate:
	PATH=$$PATH:./node_modules/.bin buf generate && pnpm buf generate
