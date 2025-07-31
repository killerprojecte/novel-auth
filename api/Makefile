.PHONY: generate

start_debug:
	cd .. && docker compose -f docker-compose.yml -f docker-compose.debug.yml up -d --build

start_release:
	cd .. && docker compose pull && docker compose up -d

stop:
	cd .. && docker compose down

integration_test: start_debug
	go clean -testcache
	go test ./... -v -p 4

generate:
	./script/generate_sql.sh
