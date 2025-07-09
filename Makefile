MIGRATIONS_DIR=./migrations
DOCKER_DB_URL=postgres://postgres:postgres@postgres:5432/go_messenger?sslmode=disable
PROTO_DIR := proto
PROTO_FILES := $(shell find $(PROTO_DIR) -name "*.proto")
GO_OUT := .
PROTOC := protoc
SHELL := /bin/bash

.PHONY: proto generate-mocks

proto:
	$(PROTOC) \
		--go_out=$(GO_OUT) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)

generate-mocks:
	@for dir in */; do \
		find "$$dir/internal" -type d -name ports | while read portsdir; do \
			echo "📦 $$dir -> найдена ports: $$portsdir"; \
			mocksdir=$$(dirname $$portsdir)/mocks; \
			mkdir -p "$$mocksdir"; \
			for src in $$portsdir/*.go; do \
				filename=$$(basename $$src .go); \
				mockgen -source=$$src -destination=$$mocksdir/$${filename}_mock.go -package=mocks; \
				echo "  ✅ мок для $$src → $${filename}_mock.go"; \
			done; \
		done; \
	done

generate-grpc-mocks:
	@for dir in proto/*/; do \
		pkg=$$(basename $$dir); \
		protofile=$$(find $$dir -name "*.proto"); \
		if [ -n "$$protofile" ]; then \
			echo "📦 Обрабатываем $$protofile"; \
			protoc --go_out=. --go-grpc_out=. --go_opt=paths=source_relative --go-grpc_opt=paths=source_relative $$protofile; \
			mockdir=proto/$$pkg/mocks; \
			mkdir -p "$$mockdir"; \
			service_name=$$(grep -oP 'service\s+\K\w+' $$protofile | head -n 1); \
			mockgen "github.com/Vovarama1992/go-ai-messenger/proto/$$pkg" "$${service_name}Client" > "$$mockdir/$$pkg""_mock.go"; \
			echo "  ✅ сгенерированы моки: $$mockdir/$$pkg""_mock.go"; \
		fi \
	done

swagger-init:
	@for dir in */; do \
		path=$$(find $$dir -type f -name routes.go | grep '/internal/.*/http/routes.go' | head -n 1); \
		if [ ! -z "$$path" ]; then \
			echo "📄 Swagger генерируется в $$dir (из $$path)"; \
			pkg_path=$$(echo $$path | sed -E 's|.*/internal/(.*)/http/routes.go|\1/http/routes.go|'); \
			cd $$dir && swag init --parseDependency --parseInternal -g internal/$$pkg_path && cd ..; \
		else \
			echo "⚠️  Не найден routes.go в $$dir/internal/**/http/"; \
		fi; \
	done
	@echo "🧹 Чистим LeftDelim / RightDelim из всех docs.go..."
	@find . -type f -name docs.go | while read file; do \
		sed -i '/LeftDelim/d' $$file; \
		sed -i '/RightDelim/d' $$file; \
		echo "✅ Пофикшен $$file"; \
	done

migrate-up:
	docker run --rm \
		--network go-ai-messenger_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DOCKER_DB_URL)" \
		up

migrate-down:
	docker run --rm \
		--network go-ai-messenger_default \
		-v $(PWD)/migrations:/migrations \
		migrate/migrate \
		-path=/migrations \
		-database "$(DOCKER_DB_URL)" \
		down

list-tests:
	find . -type f -name '*_test.go'

test:
	go test ./... -cover -v

test-integration:
	go test -tags=integration ./... -v