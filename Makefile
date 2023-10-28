grade:
	@go run cmd/grader/main.go

build-win:
	@GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o bin/grader.exe grader/main.go \
	&& upx bin/grader.exe

run-server:
	@modd -f ./.modd/server.modd.conf

db-migrate-up:
	@go run cmd/migration/migration.go

buf-generate:
	PATH=$$PATH:./node_modules/.bin buf generate && pnpm buf generate