build:
    go build

run:
    go build && ./rss-aggregator.exe

sqlc:
    sqlc generate

migration-up:
    cd sql/schema && goose postgres postgres://postgres:postgres@192.168.1.100:5432/rssagg up

migration-down:
    cd sql/schema && goose postgres postgres://postgres:postgres@192.168.1.100:5432/rssagg down
